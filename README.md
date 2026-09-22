# User Management API

RESTful API in Go for managing users, backed by MongoDB, with JWT-based authentication and a hexagonal (ports & adapters) architecture.

## Tech Stack and Lib

- Go 1.25
- MongoDB (official `mongo-driver/v2`)
- JWT (`golang-jwt/jwt/v5`), signed with HMAC (HS256)
- `go-playground/validator` for request validation
- `golang.org/x/crypto/bcrypt` for password hashing
- Standard library `net/http` (no web framework)

## Prerequest

- Go 1.25+
- Docker & Docker Compose (recommended) — or a local MongoDB instance if running the API outside a container

## Configuration

Environment variables (see `.env` for local defaults):

| Variable | Description | Default |
| --- | --- | --- |
| `APP_PORT` | HTTP port the API listens on | `8080` |
| `SECRET_KEY` | HMAC secret used to sign/verify JWTs | `secretKey` |
| `TOKEN_TTL` | JWT expiry duration (Go duration string, e.g. `24h`) | `24h` |
| `MONGO_URI` | Mongo connection string template, with `{user}`/`{password}`/`{host}`/`{port}` placeholders | `mongodb://{user}:{password}@{host}:{port}/?authSource=admin` |
| `MONGO_HOST` | Mongo host, repace into `MONGO_URI` | `localhost` |
| `MONGO_PORT` | Mongo port, repace into `MONGO_URI` | `27017` |
| `MONGO_USER` | Mongo username, repace into `MONGO_URI` | `admin` |
| `MONGO_PASSWORD` | Mongo password, repace into `MONGO_URI` | `secret123` |
| `MONGO_MINPOOL_SIZE` | Connection pool minimum size | `3` |
| `MONGO_MAXPOOL_SIZE` | Connection pool maximum size | `30` |
| `MONGO_MAX_IDLE_TIME` | Max idle time for a pooled connection (Go duration string) | `30s` |
| `MONGO_CONNECT_TIMEOUT` | Connection timeout (Go duration string) | `10s` |

The Mongo database name is fixed in code as `user_management`; the collection is `user`.

## Running with Docker Compose

Starts both MongoDB and the API:

```bash
docker compose up -d
```

## Running Locally (without Docker)

1. Start MongoDB.
```bash
cd docker/mongoDB
docker compose up -d
```
2. Ensure `.env` points `MONGO_HOST`/`MONGO_PORT` at that instance.
3. Run:

```bash
go run ./cmd
```

The server listens on `APP_PORT` (default `8080`)

## Documentation

- [API examples](documents/api-examples.md) — sample requests/responses for every endpoint
- [JWT guide](documents/jwt-guide.md) — how to get and use a token
- [Design decisions](documents/design-decisions.md) — assumptions and trade-offs made during implementation
- [Local development credentials](documents/credentials.md) — local dev credentials and where they're set