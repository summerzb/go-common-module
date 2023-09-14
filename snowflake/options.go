package snowflake

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/redis/go-redis/v9"
)

type Options struct {
	appName  string
	stepBits uint8
	nodeBits uint8
	nodeMax  int64
	epoch    int64
	gf       genFunc

	redis *redis.Client
}

type OptionFunc func(*Options)

type genFunc func() int64

// SetAppName appName 用于区别服务
func SetAppName(s string) OptionFunc {
	return func(o *Options) {
		o.appName = s
	}
}

// SetRedis 使用 redis 自增设置 node
func SetRedis(r *redis.Client) OptionFunc {
	return func(o *Options) {
		o.redis = r
	}
}

// SetEpoch 设置起始时间戳
//
//goland:noinspection ALL
func SetEpoch(n int64) OptionFunc {
	return func(o *Options) {
		if n > 0 {
			o.epoch = n
		}
	}
}

// SetNodeBits 设置数据机器位，最大 10 位(1-10)，可以部署 1024 个节点
func SetNodeBits(n uint8) OptionFunc {
	return func(o *Options) {
		if n > 0 && n <= maxNodeBits {
			o.nodeBits = n
			o.nodeMax = -1 ^ (-1 << n)
		}
	}
}

// SetStepBits 设置计数序列码，最大 12 位(1-12)，每毫秒产生 4096 个 ID
func SetStepBits(n uint8) OptionFunc {
	return func(o *Options) {
		if n > 0 && n <= maxStepBits {
			o.stepBits = n
		}
	}
}

// SetGenFunc 自定义生成节点方法
func SetGenFunc(fn genFunc) OptionFunc {
	return func(o *Options) {
		o.gf = fn
	}
}

// 默认值
func buildDefaultOptions() Options {
	return Options{
		stepBits: maxStepBits,
		nodeBits: maxNodeBits,
		nodeMax:  bits2Int64(maxNodeBits),
		epoch:    1648171610,
		gf:       randomGenFunc,
	}
}

func randomGenFunc() int64 {
	i := rand.NewSource(time.Now().UnixNano())
	return i.Int63()
}

func bits2Int64(n uint8) int64 {
	return -1 ^ (-1 << n)
}

// node 初始化
func newSnowflakeNode(id int64, o Options) (*snowflake.Node, error) {
	snowflake.Epoch = o.epoch
	snowflake.NodeBits = o.nodeBits
	snowflake.StepBits = o.stepBits
	nodeID := id % o.nodeMax

	return snowflake.NewNode(nodeID)
}
