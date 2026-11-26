# Authorizer Service

An OAuth 2.0 / OpenID Connect (OIDC) authorization server built with Go, following Clean Architecture principles. It provides user authentication, token management, email verification, password recovery, and role-based access control (RBAC) for multi-tenant applications.

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Environment Variables](#environment-variables)
  - [Database Migrations](#database-migrations)
  - [Seeding](#seeding)
  - [Running Locally](#running-locally)
- [API Endpoints](#api-endpoints)
- [Runtime Modes](#runtime-modes)
- [Docker](#docker)
- [CI/CD](#cicd)
- [Database Schema](#database-schema)
- [License](#license)

## Features

- **OAuth 2.0 Authorization Code Flow** with PKCE support
- **OpenID Connect Discovery** (`/.well-known/openid-configuration`, JWKS)
- **JWT-based Access & ID Tokens** with RSA signing (RS256)
- **User Registration** with email verification
- **Password Reset** flow via email
- **Token Exchange, Refresh & Revocation**
- **Multi-tenant RBAC** — Organizations, Roles, Permissions
- **Asynchronous Email Delivery** via RabbitMQ (publisher/consumer pattern)
- **Rate Limiting** per endpoint (Redis-backed sliding window)
- **Swagger / OpenAPI** documentation (non-production environments)
- **Health & Readiness Probes** (`/health`, `/ready`)
- **Structured JSON Logging** via `slog`
- **Database Seeding** — bootstrap admin user, application, roles & permissions
- **Docker-ready** with multi-stage build and non-root user
- **Graceful Shutdown** with signal handling

## Tech Stack

| Category | Technology |
|---|---|
| Language | Go 1.27 |
| HTTP Framework | [Gin](https://github.com/gin-gonic/gin) |
| Database (relational) | PostgreSQL + [GORM](https://gorm.io) |
| Cache / Sessions | Redis ([go-redis](https://github.com/redis/go-redis)) |
| Message Broker | RabbitMQ ([amqp091-go](https://github.com/rabbitmq/amqp091-go)) |
| Authentication | JWT ([golang-jwt](https://github.com/golang-jwt/jwt)) |
| Email | SMTP (via RabbitMQ consumer) |
| Validation | [go-playground/validator](https://github.com/go-playground/validator) |
| API Docs | [Swaggo](https://github.com/swaggo/swag) + gin-swagger |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Security | CORS, Helmet-style headers (`gin-contrib/secure`), bcrypt password hashing, SHA-256 token hashing |
| Hot Reload (dev) | [Air](https://github.com/air-verse/air) |
| Container | Docker (multi-stage Alpine) |
| CI/CD | GitHub Actions |

## Architecture

The project follows **Clean Architecture** (Hexagonal / Ports & Adapters), keeping the domain layer independent of infrastructure concerns.

```
┌──────────────────────────────────────────────────────┐
│                    Interfaces Layer                   │
│  ┌──────────────┐  ┌──────────────┐  ┌────────────┐  │
│  │  HTTP Public  │  │ HTTP Private │  │   Event    │  │
│  │  (Handlers)   │  │  (Handlers)  │  │ (Consumer) │  │
│  └──────┬───────┘  └──────┬───────┘  └─────┬──────┘  │
│         │                 │                │          │
├─────────┴─────────────────┴────────────────┴──────────┤
│                     Usecase Layer                     │
│         (Application Business Logic)                  │
├───────────────────────────────────────────────────────┤
│                     Domain Layer                      │
│  ┌────────────┐  ┌──────────────┐  ┌───────────────┐  │
│  │  Entities   │  │ Repository   │  │   Service     │  │
│  │            │  │ Interfaces   │  │  Interfaces   │  │
│  └────────────┘  └──────────────┘  └───────────────┘  │
├───────────────────────────────────────────────────────┤
│                 Infrastructure Layer                  │
│  ┌────────────┐  ┌──────────────┐  ┌───────────────┐  │
│  │  Postgres   │  │    Redis     │  │   RabbitMQ    │  │
│  │  Repos      │  │    Repos     │  │  Publisher/   │  │
│  │            │  │              │  │  Consumer     │  │
│  └────────────┘  └──────────────┘  └───────────────┘  │
│  ┌────────────┐  ┌──────────────┐  ┌───────────────┐  │
│  │   JWT /     │  │    SMTP      │  │   Logger      │  │
│  │   OIDC      │  │   Sender     │  │   (slog)      │  │
│  └────────────┘  └──────────────┘  └───────────────┘  │
└───────────────────────────────────────────────────────┘
```

**Data flow for email delivery:**

```
Register/ForgotPassword (Usecase)
  → RabbitMQ Publisher (publishes to "authorizer.email" exchange)
    → email.send queue
      → EVENT worker (RabbitMQ Consumer)
        → SMTP Sender (gmail, etc.)
          → User's inbox
```

## Project Structure

```
authorizer-service/
├── cmd/
│   ├── api/                    # Main application entrypoint
│   │   └── main.go             # Starts HTTP_PUBLIC, HTTP_PRIVATE, or EVENT worker
│   └── seed/                   # Database seeding CLI
│       └── main.go
├── internal/
│   ├── bootstrap/              # Dependency injection container
│   │   └── container.go
│   ├── config/                 # Configuration loaders (env-based)
│   │   ├── config.go           # Root config struct & Load()
│   │   ├── postgres.go
│   │   ├── redis.go
│   │   ├── rabbitmq.go
│   │   ├── smtp.go
│   │   ├── jwt.go
│   │   ├── oidc.go
│   │   ├── httppublic.go
│   │   ├── httpprivate.go
│   │   └── seed.go
│   ├── definition/
│   │   └── enum/               # Error code enumerations
│   ├── domain/                 # Core domain — no external dependencies
│   │   ├── entity/             # Domain entities (User, Organization, Role, etc.)
│   │   ├── repository/         # Repository interfaces (ports)
│   │   │   ├── user/
│   │   │   ├── application/
│   │   │   ├── organization/
│   │   │   ├── role/
│   │   │   ├── permission/
│   │   │   ├── emailverificationtoken/
│   │   │   ├── passwordresettoken/
│   │   │   └── ...
│   │   ├── service/            # Service interfaces (JWT, EmailSender, Logger)
│   │   └── apperr/             # Application error types
│   ├── infrastructure/         # External adapters (implementations)
│   │   ├── auth/               # JWT signing, JWKS, OIDC, auth code service
│   │   ├── driver/
│   │   │   ├── database/       # Postgres & Redis connection managers
│   │   │   ├── email/          # SMTP sender + RabbitMQ publisher
│   │   │   ├── rabbitmq/       # RabbitMQ connection & topology setup
│   │   │   └── logger/         # Structured JSON logger (slog)
│   │   ├── persistence/
│   │   │   ├── postgres/
│   │   │   │   ├── migrations/ # SQL migration files (golang-migrate)
│   │   │   │   ├── model/      # GORM models
│   │   │   │   └── repository/ # Postgres repository implementations
│   │   │   └── redis/
│   │   │       └── repository/ # Redis repository implementations
│   │   └── seed/               # Database seeder logic
│   ├── interfaces/             # Input/Output adapters
│   │   ├── http/
│   │   │   ├── middleware/     # RequestID, AccessLog, JWT Auth, Rate Limiter, etc.
│   │   │   ├── serializer/    # Response DTOs
│   │   │   ├── server/        # HTTP server lifecycle (graceful shutdown)
│   │   │   └── public/
│   │   │       ├── handler/v1/
│   │   │       │   ├── auth/       # Auth handlers (login, register, etc.)
│   │   │       │   ├── token/      # Token handlers (exchange, refresh, revoke)
│   │   │       │   └── discovery/  # OIDC discovery & JWKS
│   │   │       ├── router/v1/      # Route definitions
│   │   │       └── launch.go       # HTTP server bootstrap
│   │   └── event/
│   │       ├── launch.go           # Event worker bootstrap
│   │       └── email/              # RabbitMQ email consumer
│   └── usecase/                # Application use cases
│       ├── auth/               # Authorize, Login, Consent, Register, etc.
│       ├── token/              # Exchange, Refresh, Revoke
│       └── user/               # User management (admin)
├── pkg/                        # Reusable utility packages
│   ├── envutil/                # Environment variable helpers
│   ├── hash/                   # Bcrypt & SHA-256 hashing
│   ├── randutil/               # Secure random string generation
│   ├── response/               # Standardized HTTP response helpers
│   ├── rsakey/                 # RSA key loading (PEM file or string)
│   ├── stringopr/              # String operations
│   └── validation/             # Custom validator setup
├── docs/                       # Swagger generated docs
├── .air.toml                   # Air hot-reload configuration
├── .env                        # Environment variables (not committed in production)
├── Dockerfile                  # Multi-stage Docker build
├── Makefile                    # Dev & build commands
├── go.mod
└── go.sum
```

## Getting Started

### Prerequisites

- **Go** ≥ 1.27
- **PostgreSQL** ≥ 14
- **Redis** ≥ 6
- **RabbitMQ** ≥ 3.12
- **golang-migrate** CLI — `go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
- **Air** (optional, for hot reload) — `go install github.com/air-verse/air@latest`

### Environment Variables

Copy `.env.example` to `.env` (or create `.env` with the following variables):

| Variable | Description | Example |
|---|---|---|
| `POSTGRES_HOST` | PostgreSQL host | `localhost` |
| `POSTGRES_PORT` | PostgreSQL port | `5432` |
| `POSTGRES_USER` | Database user | `postgresdb` |
| `POSTGRES_PASSWORD` | Database password | `yourpassword` |
| `POSTGRES_DATABASE` | Database name | `devdb` |
| `POSTGRES_SSLMODE` | SSL mode | `disable` |
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PORT` | Redis port | `6379` |
| `RABBITMQ_HOST` | RabbitMQ host | `localhost` |
| `RABBITMQ_PORT` | RabbitMQ port | `5672` |
| `RABBITMQ_USER` | RabbitMQ username | `guest` |
| `RABBITMQ_PASSWORD` | RabbitMQ password | `guest` |
| `RABBITMQ_VHOST` | RabbitMQ virtual host | `/` |
| `SMTP_HOST` | SMTP server host | `smtp.gmail.com` |
| `SMTP_PORT` | SMTP server port | `587` |
| `SMTP_USERNAME` | SMTP auth username | `you@gmail.com` |
| `SMTP_PASSWORD` | SMTP auth password | `app-password` |
| `SMTP_FROM` | Sender email address | `you@gmail.com` |
| `SMTP_FROM_NAME` | Sender display name | `Authorizer` |
| `JWT_PRIVATE_KEY` | RSA private key (PEM string) | `-----BEGIN PRIVATE KEY-----\n...` |
| `JWT_PUBLIC_KEY` | RSA public key (PEM string) | `-----BEGIN PUBLIC KEY-----\n...` |
| `OAUTH_ISSUER_URL` | OIDC issuer URL | `http://localhost:38212` |
| `OAUTH_CODE_EXPIRY` | Authorization code TTL | `30m` |
| `OAUTH_REQUIRE_PKCE` | Require PKCE for auth code flow | `true` |
| `HTTP_PUBLIC_PORT` | Public HTTP server port | `38212` |
| `HTTP_PRIVATE_PORT` | Private HTTP server port | `38207` |
| `SEED_ADMIN_EMAIL` | Admin user email (seeder) | `admin@example.com` |
| `SEED_ADMIN_PASSWORD` | Admin user password (seeder) | `securepassword` |
| `SEED_ADMIN_NAME` | Admin user display name | `System Admin` |
| `SEED_APP_CLIENT_SECRET` | Application client secret (seeder) | `randomsecret` |

### Database Migrations

```bash
# Run all pending migrations
make migrate-up

# Roll back all migrations
make migrate-down

# Create a new migration
make migrate-create MIGRATION_NAME=create_table_xyz
```

### Seeding

Seed the database with the default application, admin user, roles, and permissions:

```bash
make seed
```

### Running Locally

The service has two runtime modes that must run simultaneously: an **HTTP server** and an **Event worker** (RabbitMQ consumer for email delivery).

```bash
# Start both HTTP server and Event worker (recommended)
make local

# Or start them individually in separate terminals:
make local-http-public   # Starts HTTP server on HTTP_PUBLIC_PORT
make local-event         # Starts RabbitMQ email consumer
```

> **Note:** If only the HTTP server is running, emails (verification, password reset) will be published to RabbitMQ but not delivered until the Event worker is started.

## API Endpoints

### OIDC Discovery

| Method | Path | Description |
|---|---|---|
| `GET` | `/.well-known/openid-configuration` | OpenID Connect discovery document |
| `GET` | `/.well-known/jwks.json` | JSON Web Key Set |

### OAuth 2.0 / Auth

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/oauth2/authorize` | — | Start authorization code flow |
| `POST` | `/oauth2/authorize/login` | — | Submit login credentials (rate limited) |
| `POST` | `/oauth2/authorize/consent` | — | Submit user consent |
| `POST` | `/oauth2/register` | — | Register a new user (rate limited) |
| `GET` | `/oauth2/verify-email` | — | Verify email via token |
| `POST` | `/oauth2/password/forgot` | — | Request password reset email (rate limited) |
| `POST` | `/oauth2/password/reset` | — | Reset password with token |
| `GET` | `/oauth2/userinfo` | Bearer JWT | Get authenticated user info |
| `POST` | `/oauth2/logout` | Bearer JWT | Logout / revoke session |

### Token

| Method | Path | Description |
|---|---|---|
| `POST` | `/oauth2/token` | Exchange authorization code for tokens (rate limited) |
| `POST` | `/oauth2/revoke` | Revoke a refresh token |

### Infrastructure

| Method | Path | Description |
|---|---|---|
| `GET` | `/health` | Liveness probe |
| `GET` | `/ready` | Readiness probe (checks Postgres & Redis) |
| `GET` | `/swagger/*` | Swagger UI (non-production only) |

## Runtime Modes

The application binary supports three runtime modes, selected via the `INTERFACE` environment variable:

| Mode | `INTERFACE` | Description |
|---|---|---|
| **HTTP Public** | `HTTP_PUBLIC` | Public-facing REST API (OAuth2, OIDC, auth endpoints) |
| **HTTP Private** | `HTTP_PRIVATE` | Internal/admin API (planned) |
| **Event Worker** | `EVENT` | RabbitMQ consumer for async tasks (email delivery) |

All modes share the same binary and configuration, differing only in which interface is activated. In production (Kubernetes), each mode runs as a separate deployment/pod.

## Docker

### Build

```bash
docker build \
  --build-arg VERSION=1.0.0 \
  --build-arg COMMIT=$(git rev-parse --short HEAD) \
  -t authorizer-service:latest .
```

### Run

```bash
# HTTP Public server
docker run -e INTERFACE=HTTP_PUBLIC --env-file .env -p 8181:8181 authorizer-service:latest

# Event worker
docker run -e INTERFACE=EVENT --env-file .env authorizer-service:latest

# Run database migrations
docker run --env-file .env --entrypoint migrate authorizer-service:latest \
  -path /app/migrations -database "$DB_URL" up

# Run database seeder
docker run --env-file .env --entrypoint /app/seed authorizer-service:latest
```

The Docker image includes:
- `/app/authorizer` — Main application binary
- `/app/seed` — Database seeder binary
- `/usr/local/bin/migrate` — golang-migrate CLI
- `/app/migrations/` — SQL migration files

## CI/CD

A GitHub Actions workflow (`.github/workflows/release.yaml`) automatically builds and pushes Docker images to DockerHub on:

- **Push to `dev` branch** — Tagged as `<version>-dev.<short-commit>`
- **Git tags** (`*.*.*`) — Tagged as the version + `latest`

Image: `docker.io/mafzaidi93/authorizer-service`

## Database Schema

The service manages the following tables (via sequential SQL migrations):

| Table | Description |
|---|---|
| `users` | User accounts (email, password hash, status) |
| `organizations` | Multi-tenant organizations |
| `organization_users` | User ↔ Organization membership |
| `applications` | OAuth2 client applications |
| `organization_applications` | Organization ↔ Application associations |
| `application_scopes` | OAuth2 scopes per application |
| `roles` | Roles (per organization or global) |
| `permissions` | Fine-grained permissions (resource + action) |
| `role_permissions` | Role ↔ Permission mapping |
| `user_roles` | User ↔ Role assignment |
| `oauth_authorization_codes` | OAuth2 authorization codes |
| `password_reset_tokens` | Password reset tokens (SHA-256 hashed) |
| `email_verification_tokens` | Email verification tokens (SHA-256 hashed) |
| `invitations` | Organization invitations |
| `audit_logs` | Audit trail for security-sensitive actions |

Additionally, **Redis** is used for:
- OAuth2 authorization codes (short-lived)
- Authorization sessions (login/consent flow state)
- OAuth2 refresh tokens
- Rate limiting counters

## License

[Apache 2.0](http://www.apache.org/licenses/LICENSE-2.0)
