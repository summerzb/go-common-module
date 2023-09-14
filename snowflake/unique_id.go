package snowflake

import (
	"errors"

	"github.com/bwmarrin/snowflake"
)

var (
	ErrEmptyUniqueID = errors.New("聚合根/实体标识为空")
)

// UniqueID 唯一标识
type UniqueID snowflake.ID

func (v UniqueID) Int64() int64 {
	return int64(v)
}

func (v UniqueID) UInt64() uint64 {
	return uint64(v)
}

func (v UniqueID) Equal(id UniqueID) bool {
	return v == id
}

func (v UniqueID) Validate() error {
	if v.IsEmpty() {
		return ErrEmptyUniqueID
	}
	return nil
}

func (v UniqueID) IsEmpty() bool {
	return v == 0
}

//goland:noinspection ALL
func ParseUniqueID(id uint64) UniqueID {
	return UniqueID(id)
}
