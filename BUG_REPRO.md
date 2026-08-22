# BUG_REPRO

## 零值路径静默放行

## Bug 是什么
`loadConfig` 缺省端口零值未兜底；`validateRunStatus` 接受空状态；规则 0704 缺省字段为零值导致校验被跳过且缺失判定反转；`ObservationRun.EnsureDefaults` 零值路径不补默认值。

## 如何触发
- 未设置 PORT 启动或新建不带 labels 的记录（可复现：`go test . -run '^TestConfigDefaultsPort$' -count=1`）。

## 错误信息
```
默认端口为 ""（期望 8080）
validateRunStatus("") 返回 nil（期望报错）
rule 0704 RequiredLabels 为 nil，校验静默放行
```
