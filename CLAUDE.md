# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Purpose

Backend for the **banking ops back-office portal** (front end is built separately in Lovable; this repo is APIs only). Scope: Thredd reconciliation (auth/financial checks vs ledger), settlement file processing, PSF/fee reporting. Deliberately split from `banking-svc` so ops tooling is isolated from the customer app. Created from `wahed-boilerplate`.

## Commands

```bash
make test        # ENV=test go test -tags="unit functional integration" ./...
make static-check
ENV=local go run main.go migratedb up && ENV=local go run main.go serve   # needs local Postgres, see README
```

## Architecture

- `cmd/serve.go` — composition root: Postgres pool → REST server (`api.NewRouter`, port 8080) + Kafka consumer (`proc/`, still the boilerplate example processor on `banking_ops_events`).
- `api/` — REST handlers on a stdlib `http.ServeMux` built from go-lib's `server.NewHealthRouter()` (so `/healthz` stays) and wrapped in `server.WithTracing`. Routes use Go 1.22 method patterns (`"GET /v1/ping"`). Version new routes under `/v1/`.
- `repository/` — **Postgres** via `pgx/v5` `pgxpool` (not Elasticsearch, unlike most services).
- `db/migrations/` — golang-migrate SQL files `NNNNNN_name.{up,down}.sql`; add the next number, never renumber. Applied by `app migratedb up` (k8s initContainer).
- `conf/` — `conf.R`; Postgres creds/host/db come from `POSTGRES_*` env vars (k8s configmap + ExternalSecret).

## Data rules

Store only what recon/settlement needs (transaction id, amount, currency, direction, fees, dates). **Do not persist cardholder PII or full raw Thredd records** without compliance sign-off.
