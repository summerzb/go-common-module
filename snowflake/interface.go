package snowflake

// NodeIface SnowFlakeNode
type NodeIface interface {
	// Init 初始化
	Init()
	// GenerateID 生成 UniqueID
	GenerateID() int64
}

// New 默认使用 redis 自增做节点
func New(opts ...OptionFunc) NodeIface {
	sf := newRedisNode(opts...)
	sf.Init()
	return sf
}

// NewCustomNode 自定义
func NewCustomNode(opts ...OptionFunc) NodeIface {
	sf := newCustomNode(opts...)
	sf.Init()
	return sf
}
