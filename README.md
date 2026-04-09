# Go Clean Architecture Boilerplate

A production-ready Go boilerplate implementing clean architecture principles with Fiber v3, PostgreSQL, Redis, JWT auth, Prometheus metrics, structured logging, and Viper-based configuration.

## Overview

This repository demonstrates a modular clean architecture design for a REST API.

- `cmd/` contains application entrypoints.
- `internal/` contains application-specific logic, split into delivery, use case, repository, entity, and model layers.
- `pkg/` contains reusable infrastructure helpers such as database setup, JWT generation, logging, metrics, validation, and error handling.
- `configs/` contains the runtime configuration schema.
- `api/openapi.yaml` serves API documentation.

### Architecture diagram (text-based)

```
Client --> Fiber HTTP Router --> HTTP Handlers
             |                 |
             v                 v
          Middleware        Use Cases
                               |
                               v
                          Domain Entities
                               |
                  +------------+-------------+
                  |                          |
             Postgres Repo              Redis Repo
                  |                          |
           SQLC / pgx / migrations      Session store
```

## Tech stack

- Go
- Fiber v3
- PostgreSQL
- Redis
- AWS S3 (aws-sdk-go-v2)
- sqlc
- pgx
- JWT
- Prometheus
- slog
- Viper

## Project structure

- `cmd/app/main.go` - application bootstrap and server startup.
- `cmd/seeder/main.go` - seed runner for test/demo data.
- `configs/` - environment-driven configuration files.
- `internal/app/` - app wiring, dependency injection, database initialization, server registration.
- `internal/delivery/http/` - HTTP delivery layer with Fiber server, router, handlers, middleware, response helpers.
- `internal/entity/` - domain entity definitions and repository contracts.
- `internal/model/` - request/response models and mappers.
- `internal/repository/` - database and cache persistence implementations.
  - `postgres/` - PostgreSQL repository implementations and transaction support.
  - `redis/` - Redis session repository.
- `internal/usecase/` - business logic layer for auth, user, and role flows.
- `pkg/storage/` - file storage implementations (AWS S3).
- `pkg/apperror/` - standardized application error package.
- `pkg/database/` - database connection, migration, and cleanup helpers.
- `pkg/jwt/` - JWT generation and parsing utilities.
- `pkg/logger/` - slog initialization.
- `pkg/metrics/` - Prometheus metric helpers.
- `pkg/validator/` - validation helpers and tests.
- `api/openapi.yaml` - OpenAPI specification and documentation source.

## Prerequisites

- Go 1.22+ (or compatible version)
- PostgreSQL
- Redis
- Docker / Docker Compose (optional for local development)
- `migrate` CLI for SQL migrations if running migrations manually
- `sqlc` if regenerating SQL client code

## Getting started

1. Clone the repository

```bash
git clone https://github.com/arielashari/go-boilerplate-clean-architecture.git
cd go-boilerplate-clean-architecture
```

2. Copy configuration

```bash
cp configs/config.example.yaml configs/config.dev.yaml
```

3. Update `configs/config.dev.yaml` to match your local Postgres and Redis settings.

4. Start local dependencies

If you have a `docker-compose.yml` configured for Postgres and Redis, you can use:

```bash
make up
```

If the repository does not include a compose file, start Postgres and Redis manually or provide equivalent hosted services.

5. Run migrations

```bash
make migrate-up
```

6. Seed demo data

```bash
make seed
```

7. Run the application

```bash
make run
```

or build and execute:

```bash
make build
./tmp/app
```

## Available make commands

- `make run` - start the app with `air` for live reload.
- `make build` - compile the app binary into `tmp/app`.
- `make up` - run `docker compose up -d` for local dependencies.
- `make down` - run `docker compose down`.
- `make migrate-up` - run database migrations upward.
- `make migrate-down` - roll back the last migration.
- `make migrate-create` - create a new SQL migration file interactively.
- `make seed` - execute the seeder at `cmd/seeder`.
- `make generate` - run `sqlc generate`.
- `make tidy` - run `go mod tidy`.
- `make test` - execute unit tests with race detection.
- `make test-cover` - execute tests and generate HTML coverage report.

## Environment configuration

The development configuration is defined in `configs/config.dev.yaml`.

```yaml
app:
  name: "myapp"
  port: 8080
  env: "dev"

database:
  postgres:
    host: "localhost"
    port: 5432
    user: "postgres"
    password: "123"
    name: "app_db"
    sslmode: "disable"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "dev-secret-not-for-production"
  access_expire_minutes: 15
  refresh_expire_minutes: 168

ratelimit:
  max_requests: 100
  expiration_seconds: 60

cors:
  allow_origins:
    - "*"
  allow_headers:
    - "Origin"
    - "Content-Type"
    - "Accept"
    - "Authorization"

s3:
  access_key_id: "your-access-key-id"
  secret_access_key: "your-secret-access-key"
  region: "us-east-1"
  bucket: "your-bucket-name"
  base_url: "https://your-bucket-name.s3.amazonaws.com"
  presign_expiry_minutes: 15
```

## API endpoints overview

### Auth

- `POST /api/v1/auth/login` - login and retrieve access + refresh tokens.
- `POST /api/v1/auth/register` - register a new user.
- `POST /api/v1/auth/logout` - invalidate the current user session. Requires Authorization header.
- `POST /api/v1/auth/refresh` - refresh an expired access token using a refresh token.

### Users (protected)

- `POST /api/v1/users` - create a new user.
- `GET /api/v1/users/` - list users with pagination and filtering.
- `GET /api/v1/users/:id` - retrieve a user by ID.
- `PATCH /api/v1/users/:id` - update an existing user.
- `DELETE /api/v1/users/:id` - delete a user.

### Roles (protected)

- `POST /api/v1/roles` - create a new role.
- `GET /api/v1/roles/` - list roles.
- `GET /api/v1/roles/:id` - retrieve a role by ID.
- `PATCH /api/v1/roles/:id` - update a role.
- `DELETE /api/v1/roles/:id` - delete a role.

### Files (protected)

- `POST /api/v1/files/upload` - upload a file to S3 storage.
  - Form fields: `entity_type`, `entity_id`, `file` (multipart)
  - Supports image/jpeg, image/png, image/webp, application/pdf
  - Max file size: 10MB
  - Returns S3 key and public URL
- `DELETE /api/v1/files?key=<s3-key>` - delete a file from S3 storage by key.
- `GET /api/v1/files/presigned?key=<s3-key>&operation=GET|PUT` - generate a presigned URL for signed access.

### Docs & metadata

- `GET /documentation` - serve API reference from `api/openapi.yaml`.
- `GET /health` - liveness check.
- `GET /ready` - readiness check.
- `GET /metrics` - Prometheus metrics endpoint.

## Architecture layers explanation

### Entity

Holds domain objects and repository contracts.
The `internal/entity` package defines domain structures and errors used across upper layers.

### Usecase

Contains business logic and orchestration.
Use cases are the application core and invoke repository contracts without depending on HTTP or transport details.

### Repository

Implements persistence for PostgreSQL and Redis.
The `postgres` repository uses `sqlc` + `pgx` to execute SQL against Postgres.
The `redis` repository stores auth session state for refresh token validation.

### Delivery

The HTTP delivery layer exposes the API through Fiber.
Handlers validate and bind requests, call use cases, and return serialized responses.

## Transaction support

Transaction support is implemented in `internal/repository/postgres/transactor.go`.
The `AuthUseCase.Register` flow wraps user creation inside `transactor.WithTx`, which begins a transaction, executes repository actions using an injected transaction context, and commits only if all operations succeed.
If the callback returns an error, the transaction rolls back automatically.

## AWS S3 File Storage

The application integrates AWS S3 for file storage through the `FileUploadUseCase` and `S3Storage` implementation.

### Configuration

Set the following environment variables or update `configs/config.dev.yaml`:

- `s3.access_key_id` - AWS access key ID (use IAM roles in production)
- `s3.secret_access_key` - AWS secret access key (use IAM roles in production)
- `s3.region` - AWS region (e.g., us-east-1)
- `s3.bucket` - S3 bucket name
- `s3.base_url` - Public base URL for S3 objects
- `s3.presign_expiry_minutes` - Expiration time for presigned URLs (default: 15)

### Features

- **Direct streaming**: Files are streamed directly to S3 without buffering in memory
- **File validation**: Enforces maximum file size (10MB) and allowed MIME types (JPEG, PNG, WebP, PDF)
- **Unique naming**: Generates UUID-based S3 keys to prevent collisions
- **Presigned URLs**: Supports time-limited URLs for both GET (download) and PUT (upload) operations
- **Error handling**: All S3 errors wrapped with context via `apperror` package

### File Upload Flow

1. Client submits multipart form to `POST /api/v1/files/upload` with:
   - `entity_type` - classification of the entity (e.g., "user", "profile")
   - `entity_id` - ID of the owning entity
   - `file` - the file to upload

2. Handler validates form fields, opens file stream
3. Usecase validates file size and MIME type
4. Storage implementation streams file directly to S3
5. Response includes S3 key and public URL

### Production Notes

- **Credentials**: Use AWS IAM roles instead of access keys in production
- **Presigned URLs**: Current implementation returns public URLs; for time-limited signatures, implement SigV4 presigning
- **Bucket policies**: Configure bucket CORS and policies to allow your domain access
- **Storage cleanup**: Consider implementing a background job to clean up orphaned files

## Error handling

This project uses `pkg/apperror` for structured application errors.

- `apperror.AppError` contains an error code, message, internal cause, and operation metadata.
- Errors are created through helper constructors such as `apperror.New(...)`.
- Use cases enrich errors with operation context and internal errors.
- The HTTP layer propagates these errors through Fiber, letting centralized middleware or handler logic format API responses consistently.

## Metrics endpoint

The app exposes Prometheus metrics at `GET /metrics`.
This endpoint is served by `promhttp.Handler()` and exposes default Go runtime, process, and custom metrics from the application if added in `pkg/metrics`.

## Health checks

- `GET /health` - returns `{ "status": "ok", "env": "<env>" }`.
- `GET /ready` - returns `{ "status": "ready" }`.

## Running tests and coverage

Run all tests:

```bash
make test
```

Generate test coverage HTML:

```bash
make test-cover
```

## Contributing

Contributions are welcome.

This repository enforces a protected `master` branch policy. All changes must be made through pull requests.

- If you have write access, create a branch from `master`, make your changes, and open a pull request.
- If you are an outside collaborator, fork the repository, push your changes to a branch in your fork, and submit a pull request from your fork into this repository.

1. Fork the repository or branch from `master`.
2. Create a feature branch.
3. Run `go test ./...` before submitting.
4. Keep implementation details in the correct layer: delivery, usecase, repository, entity.
5. Open a PR with a clear summary of the problem and solution.

> Note: This project is designed as a boilerplate. Keep the core architecture clean and avoid leaking infrastructure details into the use case layer.
