# BUG_REPRO

## 规则统计 nil map 崩溃与时间零值级联

## Bug 是什么
`opsParseStamp` 吞掉解析错误返回零值时间，`opsAge` 对非法时间戳算出离谱时长，`opsBackoff` 不封顶，`opsRuleCounts` 未初始化 map 导致 `assignment to entry in nil map` panic，`opsRuleTerminalCount` 计数反转。

## 如何触发
- 访问规则统计（可复现：`go test . -run '^TestRuleCountsNoPanic$' -count=1`）。

## 错误信息
```
panic: assignment to entry in nil map

goroutine 35 [running]:
panic({0x100f08580?, 0x100f4bce0?})
	.../runtime/panic.go:785
astronomy-observation-service.opsRuleCounts(...)
	backend/ops_rules_09.go:132
```
