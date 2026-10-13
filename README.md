![coverage](docs/coverage.svg)

# banking-ops-svc

Backend for the banking ops back-office portal (reconciliation, settlements, reporting).
Kept separate from `banking-svc` so ops tooling never touches the customer-facing banking app.

## Run locally

```bash
docker run -d --name banking-ops-pg -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=banking_ops_svc -p 5432:5432 postgres:16
ENV=local go run main.go migratedb up
ENV=local go run main.go serve

curl localhost:8080/healthz   # OK
curl localhost:8080/v1/ping   # {"status":"ok","postgres":"ok"}
```

## Commands

- `app serve` — REST API (`/healthz`, `/v1/...`) + Kafka consumer
- `app migratedb up` — applies `db/migrations` (golang-migrate); runs as an initContainer in k8s
- `app commands-watcher` — account command events (redact data, email update)
