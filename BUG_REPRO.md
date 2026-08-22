# BUG_REPRO

## 观测记录读路径污染与并发崩溃

## Bug 是什么
`OpsStore` 的 Get/List/Update/Put 直接返回或写入内部引用（Labels 浅拷贝共享），读路径 handler 给记录打 `fetchedAt`/`listedAt` 标签会污染仓储内已存记录；`OpsAudit.For` 用 `a.events[:0]` 原地压缩破坏其他记录的历史。并发读写下触发 data race，严重时进程崩溃。

## 如何触发
- 并发刷新记录详情与列表（可复现：`go test -race . -run '^TestReadPathConcurrentSafe$' -count=1`）。
- 读一次详情再列列表，列表里出现 `fetchedAt` 标签。

## 错误信息
```
WARNING: DATA RACE
Write at 0x00c000122c90 by goroutine 8:
  runtime.mapaccess2_faststr()
      .../runtime/map_faststr.go:117
  astronomy-observation-service.(*opsAPI).handleRecord()
      backend/ops_api.go:94
Previous write at 0x00c000122c90 by goroutine 10:
  runtime.mapaccess2_faststr()
      .../runtime/map_faststr.go:117
  astronomy-observation-service.(*opsAPI).handleRecord()
      backend/ops_api.go:94
```
