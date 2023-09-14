package snowflake

import (
	"context"
	"sync"

	"github.com/bwmarrin/snowflake"
)

var _ NodeIface = &redisNode{}

type redisNode struct {
	o        Options
	node     *snowflake.Node
	nodeOnce sync.Once
	nodeErr  error
}

func (s *redisNode) Init() {
	s.nodeOnce.Do(func() {
		nid := s.genNodeID()
		s.node, s.nodeErr = newSnowflakeNode(nid, s.o)
		if s.nodeErr != nil {
			panic(s.nodeErr)
		}
	})
}

func (s *redisNode) GenerateID() int64 {
	if s.node == nil {
		s.Init()
	}

	return s.node.Generate().Int64()
}

func (s *redisNode) genNodeID() int64 {
	// 按服务名自增
	key := nodeKey(s.o.appName)
	val, err := s.o.redis.Incr(context.TODO(), key).Result()
	if err != nil {
		panic(err)
	}

	return val
}

func newRedisNode(opts ...OptionFunc) *redisNode {
	options := buildDefaultOptions()

	for _, f := range opts {
		f(&options)
	}

	node := &redisNode{
		o: options,
	}

	return node
}

// nodeKey 使用 appName 生成 redis key
func nodeKey(appName string) string {
	return nodeNamespaceKey + appName
}
