# Build stage
FROM golang:1.27.1-alpine AS builder

ARG VERSION=dev
ARG COMMIT=unknown

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT}" \
    -o /app/bin/authorizer ./cmd/api

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT}" \
    -o /app/bin/seed ./cmd/seed

# Final stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata curl

# Security: Run as non-root
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder --chown=appuser:appgroup /app/bin/authorizer /app/authorizer
COPY --from=builder --chown=appuser:appgroup /app/bin/seed /app/seed

USER appuser

# Runtime configuration
ENV PORT=4000
EXPOSE $PORT

HEALTHCHECK --interval=30s --timeout=3s \
  CMD curl -f http://localhost:$PORT/health || exit 1

ENTRYPOINT ["/app/authorizer"]