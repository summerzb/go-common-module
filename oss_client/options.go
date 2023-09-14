package oss_client

// ProgressEvent defines progress event
type ProgressEvent struct {
	ConsumedBytes int64
	TotalBytes    int64
}

// ProgressListener listens progress change
type ProgressListener interface {
	ProgressChanged(event *ProgressEvent)
}

// Options oss的参数选项
type Options struct {
	IsPrivate    bool             //  是否私有
	ContentType  string           // 类型
	Process      ProgressListener // 进度
	Prefix       string           // 前缀
	Marker       string           // 标记
	Delimiter    string           // 目录过滤
	Limit        int              // 限制列表大小
	TokenExpiry  int              // token过期时间
	SessionToken string           // 临时token
}

type Option func(o *Options)

// IsPrivate 是否私有
func IsPrivate(isPrivate bool) Option {
	return func(o *Options) {
		o.IsPrivate = isPrivate
	}
}

// TokenExpiry token过期时间
func TokenExpiry(expiry int) Option {
	return func(o *Options) {
		o.TokenExpiry = expiry
	}
}

// SessionToken 临时token
func SessionToken(token string) Option {
	return func(o *Options) {
		o.SessionToken = token
	}
}

// ContentType 设置类型
func ContentType(contentType string) Option {
	return func(o *Options) {
		o.ContentType = contentType
	}
}

// Process 设置进度
func Process(listener ProgressListener) Option {
	return func(o *Options) {
		o.Process = listener
	}
}

// Prefix 设置前缀
func Prefix(prefix string) Option {
	return func(o *Options) {
		o.Prefix = prefix
	}
}

// Marker 标记
func Marker(marker string) Option {
	return func(o *Options) {
		o.Marker = marker
	}
}

// Delimiter 过滤
func Delimiter(delimiter string) Option {
	return func(o *Options) {
		o.Delimiter = delimiter
	}
}

// Limit 限定列表大小
func Limit(limit int) Option {
	return func(o *Options) {
		o.Limit = limit
	}
}

func ForOption(opts ...Option) Options {
	op := Options{}
	for _, o := range opts {
		o(&op)
	}

	return op
}
