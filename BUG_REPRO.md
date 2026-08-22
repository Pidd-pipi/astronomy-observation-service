# BUG_REPRO

## 告警推送错误链断裂与通道状态错乱

## Bug 是什么
`AlertDispatcher.allowedSeverity` 与 `Notify` 用 `%v` 包装策略与写入错误导致 `errors.Is` 失效；`AlertSink.WriteLine` 漏判已关闭状态；`Commit` 直接置位 committed 不校验关闭；`Notify` 忽略请求 ctx 取消；`Lines` 返回内部切片。

## 如何触发
- 向已关闭的告警通道发送告警或提交（可复现：`go test . -run '^TestAlertNotifyKeepsChain$' -count=1`）。
- 使用未覆盖严重级别或已取消的 ctx 调用 Notify。

## 错误信息
```
notify on closed sink 返回错误但 errors.Is(err, ErrAlertSinkClosed) 为 false（错误链断）
write to closed sink 不报错（WriteLine 漏判关闭）
commit on closed sink 不报错且 committed 被置位
notify with cancelled ctx 不报错
```
