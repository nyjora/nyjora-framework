package nfcommon

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	ERR_UNKNOW_INTERFACE_TYPE = "unknown type of interface."
	PANIC_PUSH                = "[nfcommon] nfbuf.Push marshal panic."
	PANIC_POP                 = "[nfcommon] nfbuf.Pop unmarshal panic."
)

type Nfbuf struct {
	buf *bytes.Buffer
}

func NewNFBuf() *Nfbuf {
	return &Nfbuf{
		buf: new(bytes.Buffer),
	}
}

func NewNFBufBytes(buf []byte) *Nfbuf {
	return &Nfbuf{
		buf: bytes.NewBuffer(buf),
	}
}

func (nb *Nfbuf) Push(val interface{}) *Nfbuf {
	var err error
	switch val.(type) {
	case int8, int16, int32, int64, float32, float64, uint8, uint16, uint32, uint64:
		err = binary.Write(nb.buf, binary.LittleEndian, val)
	case []byte:
		_, err = nb.buf.Write(val.([]byte))
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
		err = binary.Read(nb.buf, binary.LittleEndian, val)
	case []byte:
		_, err = nb.buf.Read(val.([]byte))
	default:
		err = errors.New("[nfbuf] Pop unknow type.")
	}
	if err != nil {
		panic(PANIC_POP + err.Error())
	}
	return nb
}

func (nb *Nfbuf) Len() int32 {
	return int32(nb.buf.Len())
}

func (nb *Nfbuf) Bytes() []byte {
	return nb.buf.Bytes()
}

func (nb *Nfbuf) UnCompress() error {
	return nil
}

func (nb *Nfbuf) Compress() error {
	return nil
}
