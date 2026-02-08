# ShubServer
GOAL : A secure, multi user backend built in Go for managing personal productivity data, including notes, todos, passwords, and bookmarks.

This will serve as the core API for multiple clients such as CLI apps and bots.

WORK IN PROGRESS

## What it does

- `GET /` serves a simple HTML page
- `GET /health` returns `{ "status": "ok" }`
- `GET /notes` list notes (in-memory)
- `GET /notes/{id}` get one note
- `POST /notes` create a note
- `PUT /notes/{id}` update a note
- `DELETE /notes/{id}` delete a note
- `GET /static/form.html` serves a basic form page

Notes are stored in memory for now (no DB yet).

## Project layout (simple view)

- `cmd/server/main.go` main entrypoint
- `internal/app` server setup + route wiring
- `internal/health` health check handler
- `internal/notes` notes handlers + in-memory repo
- `internal/config` env config (port, env)
- `internal/db` Postgres connection helper (not wired yet)
- `internal/utils` JSON helpers
- `web/index.html` homepage
- `web/static/` static files like the form (dummy pages for now)
- `Dockerfile` multi-stage build
- `docker-compose.yml` app + Postgres
- `.env` env vars for compose

## Run it locally (Go)

```bash
# from repo root

go run ./cmd/server
```

Server prints the address, default is `http://localhost:8082` if `SERVER_PORT` is not set.

## Run with Docker Compose

```bash
# build and run app + db

docker compose up --build
```

If you only want the database:

```bash
docker compose up -d db
```

Then open a psql shell:

```bash
docker exec -it shubserver-db psql -U subrat -d shubserver
```

## Env vars

From `.env` (used by compose):

- `SERVER_PORT` (default `8082` if not set)
- `POSTGRES_DB`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_HOST`
- `POSTGRES_PORT`

## Quick curl tests

```bash
curl http://localhost:8080/health
curl http://localhost:8080/notes
```

Create a note:

```bash
curl -X POST http://localhost:8080/notes \
  -H "Content-Type: application/json" \
  -d '{"title":"first","content":"hello"}'
```
