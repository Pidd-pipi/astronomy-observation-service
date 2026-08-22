# BUG_REPRO

## 导入流程 ctx 取消与 deadline 传播失效

## Bug 是什么
`PlanIngest.Ingest` 从 `context.Background()` 派生超时导致请求 ctx 丢失、取消时返回 nil 吞掉错误；`IngestStore.Put/Get` 不检查 ctx.Done，取消后的读写仍会完成。

## 如何触发
- 已取消或已过期的请求 ctx 调用 Ingest（可复现：`go test . -run '^TestIngestDerivesFromRequest$' -count=1`）。
- 取消 ctx 后直接读写 IngestStore。

## 错误信息
```
ingest with cancelled request ctx 不返回错误（deadline 未从请求派生）
expired deadline 时 Ingest 返回 nil error（取消被吞）
cancelled Put/Get 仍然写入/读出数据
```
