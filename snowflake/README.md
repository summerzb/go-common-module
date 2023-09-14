# Snowflake 雪花算法 ID 生成器

### interface
```
// SnowFlakeNode
type NodeIface interface {
	// 初始化
	Init()
	// 生成 UniqueID
	GenerateID() UniqueID
}
```
### 初始化
```go
sf := snowflake.New(
        snowflake.SetAppName(APP_NAME),
        snowflake.SetRedis(redis),
      )

id := sf.GenerateID()
```

### Options 可选配置项
- `SetAppName` appName 用于区别服务
- `SetEpoch` 设置起始时间戳
- `SetNodeBits` 设置数据机器位，最大 10 位(1-10)，可以部署 1024 个节点
- `SetStepBits` 设置计数序列码，最大 12 位(1-12)，每毫秒产生 4096 个 ID
- `SetRedis`  redis 实例
- `SetGenFunc` 自定义

```go
sf := snowflake.New(
        snowflake.SetAppName(APP_NAME),
        snowflake.SetRedis(redis),
        snowflake.SetNodeBits(nodeBits),
        snowflake.SetStepBits(stepBits),
        snowflake.SetEpoch(epoch),
      )

id := sf.GenerateID()
```


- 使用
```go
id := vars.Snowflake.GenerateID()
```
