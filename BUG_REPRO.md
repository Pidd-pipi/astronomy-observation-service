# BUG_REPRO

## 夜间归档 defer 吞错与日志状态错乱

## Bug 是什么
`ArchiveNight` 命名返回值被 defer Close 覆盖吞错；`NightLog.Append` 漏判已关闭状态；`Commit` 先置 committed 再判关闭；空批次漏提交；`Entries` 返回内部切片。

## 如何触发
- 归档批次混入 blocked 记录（可复现：`go test . -run '^TestArchiveNightKeepsError$' -count=1`）。
- 关闭后继续写日志（可复现：`go test . -run '^TestLogAppendClosedDenied$' -count=1`）。

## 错误信息
```
archive with blocked entry 返回 nil error（错误被吞）
night log 关闭后仍可 Append / Commit，且 committed 状态提前置位
```
