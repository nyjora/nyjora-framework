package nfcommon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"sync"
)

const (
	ERR_UNKNOW_INTERFACE_TYPE = "unknown type of interface."
	PANIC_PUSH                = "[nfcommon] nfbuf.Push marshal panic."
	PANIC_POP                 = "[nfcommon] nfbuf.Pop unmarshal panic."
)

var buffPool = &sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type Nfbuf struct {
	*bytes.Buffer
}

func NewNFBuf() *Nfbuf {
	return &Nfbuf{
		buffPool.Get().(*bytes.Buffer),
	}
}

func NewNFBufNoPool() *Nfbuf {
	return &Nfbuf{
		new(bytes.Buffer),
	}
}

func FreeNFBuf(nf *Nfbuf) {
	nf.Reset()
	buffPool.Put(nf)
}

func (nb *Nfbuf) Push(val interface{}) *Nfbuf {
	var err error
	switch val.(type) {
	case int8, int16, int32, int64, float32, float64, uint8, uint16, uint32, uint64:
		err = binary.Write(nb, binary.LittleEndian, val)
	case []byte:
		_, err = nb.Write(val.([]byte))
	default:
		err = errors.New(ERR_UNKNOW_INTERFACE_TYPE)
	}
	if err != nil {
		panic(PANIC_PUSH + err.Error())
	}
	return nb
}

func (nb *Nfbuf) Pop(val interface{}) *Nfbuf {
	var err error
	switch val.(type) {
	case *int8, *int16, *int32, *int64, *float32, *float64, *uint8, *uint16, *uint32, *uint64:
		err = binary.Read(nb, binary.LittleEndian, val)
	case []byte:
		_, err = nb.Read(val.([]byte))
	default:
		err = errors.New("[nfbuf] Pop unknow type.")
	}
	if err != nil {
		panic(PANIC_POP + err.Error())
	}
	return nb
}
