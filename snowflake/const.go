package snowflake

const (
	// 机器数字位
	maxNodeBits uint8 = 10
	// 计数序列位
	maxStepBits uint8 = 12
	// redis key
	nodeNamespaceKey = "snowflakenodekey:"
)
