# BUG_REPRO

## 状态机转换表缺边与历史污染

## Bug 是什么
状态机转换表缺了 paused→active 的边导致恢复被拒；`Move` 对同状态 no-op 也写历史；`Last` 空历史越界；`opsStatusValid` 把 paused 判成非法；规则 0303 终端标记错位。

## 如何触发
- 暂停的记录点恢复（可复现：`go test . -run '^TestStatePausedResumes$' -count=1`）。
- 空状态机调用 Last（可复现：`go test . -run '^TestStateLastEmptySafe$' -count=1`）。

## 错误信息
```
panic: runtime error: index out of range [0] with length 0
astronomy-observation-service.(*OpsStateMachine).Last()
	backend/ops_state.go:66
```
