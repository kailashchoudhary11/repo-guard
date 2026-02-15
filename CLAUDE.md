# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Repo-Guard is a GitHub App that automatically detects duplicate issues using AI-powered semantic similarity. When a new issue is opened, it compares it against existing issues via a Python ML service (sentence-transformers with BAAI/bge-m3 embeddings), then comments with links to similar issues and optionally closes the duplicate.

## Development Commands

**Live reload development (Go + Templ + Tailwind):**
```bash
air
```
This runs `templ generate`, compiles Tailwind CSS, and builds/watches the Go app.

**Manual build:**
```bash
templ generate && ./tailwindcss -i ./static/css/input.css -o ./static/css/styles.css && go build -o ./tmp/main .
```

**Run services (PostgreSQL + Python ML service):**
```bash
docker-compose up
```

**Generate database code after changing SQL queries/schema:**
```bash
sqlc generate
```

**Run tests:**
```bash
go test ./...
```

## Architecture

```
GitHub (webhooks) → Go App (:8000) → Python ML Service (:8233)
                        ↕
                   PostgreSQL
```

**Go app (main service):** Receives GitHub webhooks, orchestrates duplicate detection, interacts with GitHub API. Uses goroutines with channels to compare issues concurrently.

**Python service (`python-service/`):** Flask app using sentence-transformers for semantic similarity. Endpoint `POST /compare_issues` returns similarity scores (≥0.85 = duplicate, 0.65-0.85 = possible).

**Key layers:**
- `handlers/` — HTTP handlers. `webhook.go` contains the core duplicate detection flow.
- `services/` — GitHub API interactions (fetching issues, posting comments, closing issues).
- `models/` — Data structures for GitHub webhook payloads, issues, repositories, installations.
- `initializers/` — App startup: env loading, GitHub client, OAuth2 config, DB client.
- `helpers/` — JWT generation for GitHub App authentication.
- `db/` — PostgreSQL via SQLC. Schema in `db/schema.sql`, queries in `db/queries/`, generated code in `db/generated/`.
- `templates/` — Templ components for the web UI (landing page, setup page).

**Code generation:** Templ files (`.templ`) generate `_templ.go` files. SQLC generates type-safe Go code from SQL in `db/generated/`. Neither should be hand-edited.

## Environment Variables

Copy `.env.example` and fill in: `AI_MODEL_URL`, `APP_NAME`, `APP_ID`, `CLIENT_ID`, `CLIENT_SECRET`, `PRIVATE_KEY`. Database connection requires `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_HOST`.

## Docker Services

| Service | Port (host) | Port (container) |
|---------|-------------|-------------------|
| repoguard (Go) | 8344 | 8000 |
| repo-model (Python) | 8233 | 8233 |
| postgres | 8345 | 5432 |
