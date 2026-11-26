# ---------- Build stage ----------
FROM golang:1.27.1-alpine AS builder

ARG VERSION=dev
ARG COMMIT=unknown

WORKDIR /src

# Download dependencies (cached layer)
COPY go.mod go.sum ./
RUN go mod download

# golang-migrate CLI with only the postgres driver, pinned. Installed to /go/bin.
# Built from source instead of copied from the migrate/migrate image, whose
# binary path has changed between releases.
RUN CGO_ENABLED=0 go install -tags 'postgres' -ldflags="-s -w" \
    github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.1

# Build both application binaries
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath \
    -ldflags="-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT}" \
    -o /out/authorizer ./cmd/api \
 && CGO_ENABLED=0 GOOS=linux go build -trimpath \
    -ldflags="-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT}" \
    -o /out/seed ./cmd/seed

# ---------- Final stage ----------
FROM alpine:3.22

# Fixed numeric UID/GID so Kubernetes can verify runAsNonRoot (runAsUser 10001).
RUN apk add --no-cache ca-certificates tzdata \
 && addgroup -S -g 10001 app \
 && adduser -S -D -H -u 10001 -G app -s /sbin/nologin app

# Binaries and SQL files stay owned by root, so the app user cannot modify them.
COPY --from=builder /out/authorizer   /app/authorizer
COPY --from=builder /out/seed         /app/seed
COPY --from=builder /go/bin/migrate   /usr/local/bin/migrate
COPY internal/infrastructure/persistence/postgres/migrations /app/migrations

USER 10001:10001

# INTERFACE selects the role of /app/authorizer:
#   HTTP_PUBLIC  listens on HTTP_PUBLIC_PORT  (default 8181)
#   HTTP_PRIVATE listens on HTTP_PRIVATE_PORT (default 8080)
#   EVENT        RabbitMQ consumer, no port
ENV HTTP_PUBLIC_PORT=8181
EXPOSE 8181 8080

# For `docker run` with INTERFACE=HTTP_PUBLIC only. Kubernetes ignores this and
# uses the probes from the chart.
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s \
  CMD wget -qO- "http://127.0.0.1:${HTTP_PUBLIC_PORT}/health" >/dev/null || exit 1

# Other entrypoints in this image: /app/seed (seed Job) and `migrate` (migrate Job).
ENTRYPOINT ["/app/authorizer"]
