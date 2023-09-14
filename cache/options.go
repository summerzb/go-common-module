package cache

type Options struct {
	Endpoint string
	Password string
	Db       int
	PoolSize int
	MinIdle  int
}

type Option func(o *Options)

func WithOptions(op Options) Option {
	return func(o *Options) {
		o.Endpoint = op.Endpoint
		o.Password = op.Password
		o.Db = op.Db
		o.PoolSize = op.PoolSize
		o.MinIdle = op.MinIdle
	}
}

func WithEndpoint(endpoint string) Option {
	return func(o *Options) {
		o.Endpoint = endpoint
	}
}

func WithPwd(pwd string) Option {
	return func(o *Options) {
		o.Password = pwd
	}
}

func WithDB(db int) Option {
	return func(o *Options) {
		o.Db = db
	}
}

func WithPoolSize(size int) Option {
	return func(o *Options) {
		o.PoolSize = size
	}
}

func WithMinIdle(idle int) Option {
	return func(o *Options) {
		o.MinIdle = idle
	}
}
