# Local Development Credentials

These are **local/demo credentials only** — used for running this project on your own machine. They are not production secrets, and should never be reused in a real deployment.

## Where they come from

| Credential | Used for | Set in |
| --- | --- | --- |
| `SECRET_KEY` | Signs/verifies JWTs (HMAC/HS256) | `.env` (for `go run ./cmd`), or `docker-compose.yaml` (for `docker compose up`) |
| `MONGO_USER` / `MONGO_PASSWORD` | MongoDB root login | `.env` and `docker-compose.yaml`'s `mongo` service (`MONGO_INITDB_ROOT_USERNAME`/`MONGO_INITDB_ROOT_PASSWORD`) |

## Default values in this repo

| Variable | Value | Where |
| --- | --- | --- |
| `SECRET_KEY` | `secretKey` | `.env` |
| `SECRET_KEY` | `4b949ea1-ac99-4eb4-b158-f7ae46d6fc89` | `docker-compose.yaml` (`api` service) |
| `MONGO_USER` | `admin` | `.env` / `docker-compose.yaml` |
| `MONGO_PASSWORD` | `secret123` | `.env` / `docker-compose.yaml` |

Note: `SECRET_KEY` differs between `.env` and `docker-compose.yaml` — a token issued by one won't validate against the other. Pick one path (either run everything via `docker compose up -d`, or run MongoDB in Docker and the API locally with `go run ./cmd` using `.env`) and stick to it, or align the two values if you switch between them.

## For a real deployment

Set `SECRET_KEY` to a long random value (e.g. `openssl rand -hex 32`), and set real MongoDB credentials — don't reuse anything from `.env` or `docker-compose.yaml` in this repo. See [.env.example](../.env.example) for the full list of variables to configure.
