package logger

type Config struct {
	Level       uint8  // 日志级别
	Format      uint8  // 日志格式 1: json，默认格式 2: 文件格式
	Director    string // 日志目录 默认是当前目录下的logs
	RotateType  uint8  // 1：按照日期切分, 默认 2：按照大小切分
	MaxSize     uint32 // 在进行切割之前，日志文件的最大大小（以MB为单位）
	MaxAge      uint32 // 日志最大保留时间 单位：天
	MaxBackups  uint32 // 保留旧文件的最大个数
	Compress    bool   // 是否压缩
	ShowLine    bool   // 是否在日志中输出源码所在的行
	StdoutPrint bool   // 是否输出到控制台
	AdaptType   uint8  // 适配器类型 1: zap(默认类型) 2: zerolog
}
