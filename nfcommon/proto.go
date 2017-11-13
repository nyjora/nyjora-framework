package nfcommon

import (
	"sync"
)

var protocolPool = &sync.Pool{
	New: func() interface{} {
		return &Protocol{}
	},
}

type Protocol struct {
	Id       uint32
	FromType uint32
	FromId   uint32
	ToType   uint32
	ToId     uint32
	Data     []byte
}

func NewProto() *Protocol {
	proto := protocolPool.Get().(*Protocol)
	return proto
}

func FreeProto(proto *Protocol) {
	protocolPool.Put(proto)
}
