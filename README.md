# Forgotten API

A Twitter-like social media API built with Go. This is an ongoing project; the README and features will evolve over time.

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/b107d45cd0f14123ae60578870cf9a28)](https://app.codacy.com/gh/nevzattalhaozcan/forgotten/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Coverage Badge](https://app.codacy.com/project/badge/Coverage/b107d45cd0f14123ae60578870cf9a28)](https://app.codacy.com/gh/nevzattalhaozcan/forgotten/dashboard?utm_source=github.com&utm_medium=referral&utm_content=nevzattalhaozcan/forgotten&utm_campaign=Badge_Coverage)
[![CI](https://github.com/nevzattalhaozcan/forgotten/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/nevzattalhaozcan/forgotten/actions/workflows/ci.yml)


## Overview

- REST API with JWT auth, user registration/login
- Built with Gin and Gorm (PostgreSQL)
- Structured logging, request metrics, and CORS enabled
- Swagger/OpenAPI documentation (generated with swag)
- Dockerized for local/dev; CI/CD via GitHub Actions

## Tech stack

- Go 1.23
- Gin (HTTP), Gorm (ORM, Postgres)
- Swagger (swaggo), Prometheus metrics
- Docker + docker-compose
- GitHub Actions (CI + Deploy)

## Getting started

1) Prerequisites
- Go 1.23+
- Docker (optional, for containerized dev)
- PostgreSQL (if running locally without Docker)

2) Clone and enter the project
```sh
git clone https://github.com/nevzattalhaozcan/forgotten.git
cd forgotten
```

3) Configure environment
- Copy `.env.example` to `.env` and adjust as needed:
  - `DATABASE_URL` (e.g., postgres://user:password@localhost:5432/forgotten_db?sslmode=disable)
- Seed data lives at `data/seed_users.json`.

### Run locally (without Docker)

```sh
# run the server
go run ./cmd/server
# server defaults to http://localhost:8080
```

### Run with Docker Compose

```sh
cd docker
docker-compose up -d
# API available at http://localhost:8080
```

## API docs (Swagger)

- The top-level annotations live in `cmd/server/main.go`.
- (If needed) generate/update docs:
```sh
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/server/main.go -o ./docs
```
- Swagger UI can be served at `/swagger` once wired in the router.
  - If not yet enabled, add a route that uses `gin-swagger` and import the generated `docs` package.

Base path: `/api/v1`

## Testing

- Unit tests:
```sh
go test -v -race ./...
```
- Coverage:
```sh
go test -v -race -coverprofile=coverage.out ./...
```

Note: Most tests use an in-memory SQLite DB via `pkg/testutil`; add Postgres-backed integration tests as needed.

## Metrics

- Prometheus-compatible metrics middleware is included. Expose your metrics endpoint (e.g., `/metrics`) as you wire up monitoring.

## Project structure (high level)

- `cmd/server` – application entrypoint
- `internal/` – config, database, handlers, repository, services, middleware, models
- `pkg/` – shared libs (logger, metrics, utils, test utilities)
- `docs/` – generated Swagger docs
- `docker/` – Dockerfile, docker-compose, Prometheus config
- `data/` – seed fixtures

## CI/CD

- `.github/workflows/ci.yml` runs tests, lint, and build on PRs/commits.
- `.github/workflows/deploy.yml` builds and pushes a Docker image and deploys over SSH on main (requires repo secrets).

## Roadmap

- More endpoints (profiles, posts, timelines)
- RBAC/permissions
- Full integration test suite
- Production-ready observability and hardening