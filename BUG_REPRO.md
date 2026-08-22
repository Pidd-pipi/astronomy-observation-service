# BUG_REPRO

## 批量归档并发生命周期错位

## Bug 是什么
`OpsService.ArchiveBatch` 的 wg.Add 写在 goroutine 内、错误分支漏发结果、worker 内 `defer close(results)` 与收尾 close 双重关闭，多 worker 场景 panic；结果计数也错位。

## 如何触发
- 批量归档多条记录（可复现：`go test -race . -run '^TestArchiveBatchNoPanicMany$' -count=1`）。
- 批量归档混入不存在的编号，接口不返回。

## 错误信息
```
panic: send on closed channel

goroutine 18 [running]:
astronomy-observation-service.(*OpsService).ArchiveBatch.func2()
	backend/ops_batch.go:41
```
