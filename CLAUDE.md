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

- `cmd/serve.go` — composition root: Postgres pool → auth (Supabase JWKS) + GraphQL schema → HTTP server (`api.NewRouter`, port 8080) + Kafka consumer (`proc/`, still the boilerplate example processor on `banking_ops_events`).
- `api/` — the portal API is **GraphQL only**: `GET|POST /graphql` (graphql-go, same as dashboard-svc; unversioned path; no GraphiQL, since the route is behind auth and a browser tab sends no token). It sits on a stdlib `http.ServeMux` built from go-lib's `server.NewHealthRouter()` (so `/healthz` stays), wrapped in CORS (`api/cors.go`) and `server.WithTracing`. REST is only for things GraphQL can't serve: probes (`/healthz`, `GET /v1/ping`) and, later, provider webhooks. Don't add REST twins of GraphQL queries.
- `schema/` — the GraphQL schema (`schema.New`, starter query `currentUser`); add queries to the `Query` fields, resolvers read the caller via `auth.CallerFrom`.
- `auth/` — the portal calls this service directly with the user's Supabase access token (Google SAML SSO). `auth.Authenticator` verifies it with go-lib `jwt` against the Supabase JWKS (ES256 — go-lib only allows RS256/EdDSA by default, so it is overridden — plus `iss`, `aud`, `exp`), then requires an `amr` entry `sso/saml` from our SAML provider and an `@wahed.com` email; the caller goes on the context (`auth.CallerFrom`). `/healthz` and `/v1/ping` stay open; wrap every other route in `Require`. CORS allows only `auth.allowedOrigins`. Never trust `user_metadata` claims — users can edit them.
- `repository/` — **Supabase Postgres** via `pgx/v5` `pgxpool` (not Elasticsearch, unlike most services). Supabase has no IAM auth: username/password come from the k8s secret synced by ExternalSecret from AWS Secrets Manager.
- `db/migrations/` — golang-migrate SQL files `NNNNNN_name.{up,down}.sql`; add the next number, never renumber. Applied by `app migratedb up` (k8s initContainer).
- `conf/` — `conf.R`; Postgres host/db come from `POSTGRES_*` env vars (k8s configmap via `envFrom`), username/password from the secret; the `auth` section (JWKS URL, issuer, SAML provider id, allowed origins) is set per env.

## Data rules

Store only what recon/settlement needs (transaction id, amount, currency, direction, fees, dates). **Do not persist cardholder PII or full raw Thredd records** without compliance sign-off.
