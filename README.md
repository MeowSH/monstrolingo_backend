# Monstrolingo Backend

Monstrolingo Backend is a Go + Encore API for multilingual Monster Hunter Wilds
catalog data.

It is designed to provide stable category endpoints, dynamic language options
for the frontend, and a clean deployment/setup workflow.

## Global Features

- Category APIs for:
  - `items`, `weapons`, `armor`, `skills`, `decorations`, `charms`,
    `food-skills`, `kinsects`
- Bilingual table view per category (`source_lang` + `target_lang`)
- Detailed entity view per category (`external_key` + `target_lang`)
- Dynamic language list endpoint for frontend enums (`GET /languages`)
- Link build translation endpoint (`POST /linkbuild/translate`) with:
  - strict simulator URL validation
  - automatic source language detection
  - best-effort translation with partial fallback
  - translation-only response for URL-provided entries (`skills`, set-skill effects, jewel-related entries)
- Standardized API error model (`invalid_argument`, `not_found`, `internal`)
- Health check endpoint (`GET /health`)

## High-Level Architecture

- Top-level Encore services:
  - `items`, `weapons`, `armor`, `skills`, `decorations`, `charms`,
    `foodskills`, `kinsects`, `languages`, `health`, `game`, `linkbuild`
- Shared read layer:
  - `internal/catalogcore`
- Shared sim-build translation layer:
  - `internal/simbuildcore`
- Database:
  - PostgreSQL with GORM
- Schema migrations:
  - Atlas (`db/migrations`)

## Requirements

- Go `1.24+`
- Encore CLI
- Docker + Docker Compose
- Atlas CLI

## Local Setup

### 1) Configure environment

```bash
cp .env.example .env
```

The backend auto-loads `.env` in local development, so manual `export` is not
required before `encore run`.

Expected variables:

- `POSTGRES_DB`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_PORT`
- `POSTGRES_HOST` (optional, defaults to `localhost`)
- `POSTGRES_SSLMODE` (optional, defaults to `disable`)

### 2) Start PostgreSQL

```bash
docker compose up -d postgres
```

### 3) Apply migrations

```bash
atlas migrate apply --env local
```

### 4) Run the backend

```bash
encore run
```

## Quick Verification

```bash
curl -sS "http://127.0.0.1:4000/health"
curl -sS "http://127.0.0.1:4000/languages"
curl -sS "http://127.0.0.1:4000/items/table?source_lang=en&target_lang=fr&page=1&limit=5"
curl -sS "http://127.0.0.1:4000/items/detail/potion?target_lang=fr"
curl -sS -X POST "http://127.0.0.1:4000/linkbuild/translate" \
  -H "content-type: application/json" \
  -d '{"url":"https://simulator.example/sim/#skills=Attack%20Boost%20Lv2","target_lang":"fr"}'
```

## API Contract (Frontend Summary)

### Table endpoints

Pattern:

`GET /<category>/table?source_lang=<code>&target_lang=<code>&page=<n>&limit=<n>`

Response:

- `items[]` with `external_key`, `source`, `target`
- `pagination` with `page`, `limit`, `total`, `total_pages`, `has_next`

### Detail endpoints

Pattern:

`GET /<category>/detail/:external_key?target_lang=<code>`

Response:

- `data` object containing canonical fields plus target translation

### Languages endpoint

`GET /languages`

Response:

- `languages[]` with `code` and `label`

### Sim build translation endpoint

`POST /linkbuild/translate`

Request body:

- `url` (required): must target a supported simulator URL (`/sim/` path)
- `target_lang` (required): language code from `GET /languages`

Response (high-level):

- `source_lang_detected`
- `target_lang`
- `translation_mode` (`full` or `partial`)
- `skills_original[]`
- `skills_translated[]` (`translated: true|false`)
- `unmatched_elements[]`

## Integration Conventions

- Language params are language codes (`en`, `fr`, etc.), not UUIDs
- Pagination defaults:
  - `page = 1`
  - `limit = 25`
  - max `limit = 100`
- Detail translation falls back to English when target data is missing
- Sim build translation keeps original text when elements are not translatable
- Standard error codes:
  - `invalid_argument` (400)
  - `not_found` (404)
  - `internal` (500)

## Useful Commands

```bash
go test ./...
atlas migrate status --env local
atlas migrate validate --env local
```
