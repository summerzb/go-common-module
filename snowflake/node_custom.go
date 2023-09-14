package snowflake

import (
	"sync"

	"github.com/bwmarrin/snowflake"
)

var _ NodeIface = &customNode{}

type customNode struct {
	o        Options
	node     *snowflake.Node
	nodeOnce sync.Once
	nodeErr  error
}

func (s *customNode) Init() {
	s.nodeOnce.Do(func() {
		val := s.o.gf()
		s.node, s.nodeErr = newSnowflakeNode(val, s.o)
		if s.nodeErr != nil {
			panic(s.nodeErr)
		}
	})
}

func (s *customNode) GenerateID() int64 {
	if s.node == nil {
		s.Init()
	}

	return s.node.Generate().Int64()
}

func newCustomNode(opts ...OptionFunc) *customNode {
	options := buildDefaultOptions()

	for _, f := range opts {
		f(&options)
	}

	return &customNode{
		o: options,
	}
}
