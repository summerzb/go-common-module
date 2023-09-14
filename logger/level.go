package logger

const (
	// DebugLevel 调试
	DebugLevel = iota
	// InfoLevel 信息
	InfoLevel
	// WarnLevel 警告
	WarnLevel
	// ErrorLevel 错误
	ErrorLevel
	// PanicLevel panic
	PanicLevel
	// FatalLevel 推出
	FatalLevel
)

const (
	JsonFormat = 1 // json格式
	TextFormat = 2 // txt文本
)
