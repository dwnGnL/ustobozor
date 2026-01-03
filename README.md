# Bozor API

Backend for a job marketplace where users can act as both customers and workers.

## Requirements
- Go
- PostgreSQL

## Setup
1. Copy `.env.example` to `.env` and update values.
2. Apply migrations in `migrations/` to your PostgreSQL database.
3. Run the API:

```bash
go run ./cmd/api
```

## Docker
Build and run everything without local Go/PostgreSQL:

```bash
docker compose up --build
```

To reset state:

```bash
docker compose down -v
```

Notes:
- Migrations from `migrations/` run on first DB init via `docker-entrypoint-initdb.d`.
- If you change the Go version in `go.mod`, update the `FROM golang:...` line in `Dockerfile`.

## GraphQL
- Playground: `http://localhost:8080/`
- Endpoint: `http://localhost:8080/graphql`

## Uploads
Uploaded photos are stored in `UPLOAD_DIR` and served at `/uploads/`.
