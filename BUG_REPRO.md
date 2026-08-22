# BUG_REPRO

## 错误链断链导致冲突误判 500

## Bug 是什么
`OpsError.Unwrap` 断链（errors.Is 失效）、`opsCode` 丢失包装分类、错误文本丢 cause、`opsHTTPError` 缺少 conflict 分支，重复创建记录返回 500 而不是 409。

## 如何触发
- 重复提交同编号运营记录（可复现：`go test . -run '^TestCreateDuplicate409$' -count=1`）。

## 错误信息
```
HTTP/1.1 500 Internal Server Error
{"error":"internal: create: store.put: operations revision conflict"}
```
