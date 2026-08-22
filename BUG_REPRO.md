# BUG_REPRO

## 请求上下文 deadline 与标识丢失

## Bug 是什么
`requestTimeoutMiddleware` 没有把带超时的 ctx 传给下游且中间件链漏挂，`requestIDMiddleware` 未把请求标识写入 ctx，`healthHandler` 用 `context.Background()` 判断 deadline 与请求标识，健康检查字段失真。

## 如何触发
- 访问 `GET /healthz`（可复现：`go test . -run '^TestHealthDeadlineTrue$' -count=1`）。

## 错误信息
```
GET /healthz -> {"deadline":"false","requestId":""}
期望 {"deadline":"true","requestId":"req-test-123"}
```
