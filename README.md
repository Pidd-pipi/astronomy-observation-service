# Astronomy Observation Service

Coordinates observation runs, instruments, targets, and data-quality stages for an observatory. The service also manages an operations domain for night-plan records: status transitions, rules, audit events, snapshots and batch archiving.

The `PORT` environment variable selects the HTTP port and defaults to `8080`; run with `go run .`.

## Endpoints

- `GET /healthz` — health check (reports service status, request deadline and request id).
- `GET /api/runs`, `POST /api/runs/{id}/status` — observation runs collection and status change.
- `GET /api/ops/records` — list/search operations records (`subject`, `status`, `priority`, `owner`, `page`, `pageSize`).
- `POST /api/ops/records` — create an operations record.
- `GET /api/ops/records/{id}` — fetch one record.
- `POST /api/ops/records/{id}/transition` — change record status (`{"status": "...", "expected": 0, "actor": "..."}`).
- `GET /api/ops/records/{id}/audit` — audit events for one record.
- `GET /api/ops/snapshot` — status/priority statistics.
- `GET /api/ops/rules` — operations rule set.
- `POST /api/ops/batch/archive` — batch archive (`{"ids": [...], "actor": "..."}`).

## Verification

- `go build ./...` and `go test ./...` pass in `backend/`.
- Runtime smoke: `PORT=18184 go run .`; `GET /healthz` returns HTTP 200.

## Engineering Notes

请求保留请求标识并经过恢复与超时保护；ops 域按领域模型、校验、状态转换、并发安全存储、审计事件和规则集分层。状态写入使用版本校验，错误通过可识别的领域错误返回。

## Enterprise Layout

```text
.
├── backend/                 # Go module, source code, web assets, Dockerfile
├── database/                # Database extension documentation
├── output/                  # Verification record
└── runtime_smoke.json       # Startup contract
```

Health check: `GET /healthz`. API endpoints: `GET /api/runs` and `POST /api/runs/{id}/status`.
