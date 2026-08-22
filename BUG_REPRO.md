# BUG_REPRO

## 列表筛选切片共享与分页越界

## Bug 是什么
`filterOpsRecords` 用 `items[:0]` 原地压缩共享底层数组，`Search` 总数取筛选前长度且不过滤，`opsClonePage` 直接返回内部切片，`opsBounds` 不收敛导致分页越界。

## 如何触发
- 按状态筛选后读原列表（可复现：`go test . -run '^TestFilterRecordsNoAlias$' -count=1`）。
- 大 pageSize 翻页（可复现：`go test . -run '^TestBoundsClamped$' -count=1`）。

## 错误信息
```
panic: runtime error: slice bounds out of range [:100] with capacity 4
astronomy-observation-service.(*OpsService).Search()
	backend/ops_service.go:57
```
