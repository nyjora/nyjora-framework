package nfcommon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	ERR_UNKNOW_INTERFACE_TYPE = "unknown type of interface."
	PANIC_PUSH                = "[nfcommon] nfbuf.Push marshal panic."
	PANIC_POP                 = "[nfcommon] nfbuf.Pop unmarshal panic."
)

type nfbuf struct {
	buf *bytes.Buffer
}

func NewNFBuf() *nfbuf {
	return &nfbuf{
		buf: new(bytes.Buffer),
	}
}

func (nb *nfbuf) Push(val interface{}) *nfbuf {
	var err error
	switch val.(type) {
	case int:
		err = binary.Write(nb.buf, binary.LittleEndian, val)
	case int8:
		err = binary.Write(nb.buf, binary.LittleEndian, val)
	case int16:
		err = binary.Write(nb.buf, binary.LittleEndian, val)
	case int32:
		err = binary.Write(nb.buf, binary.LittleEndian, val)
	case int64:
		err = binary.Write(nb.buf, binary.LittleEndian, val)
	case float32:
		err = binary.Write(nb.buf, binary.LittleEndian, val)
	case float64:
		err = binary.Write(nb.buf, binary.LittleEndian, val)
	case uint:
		err = binary.Write(nb.buf, binary.LittleEndian, val)
	case uint8:
		err = binary.Write(nb.buf, binary.LittleEndian, val)
	case uint16:
		err = binary.Write(nb.buf, binary.LittleEndian, val)
	case uint32:
		err = binary.Write(nb.buf, binary.LittleEndian, val)
	case uint64:
		err = binary.Write(nb.buf, binary.LittleEndian, val)
	case []byte:
		fmt.Println("write []byte")
		_, err = nb.buf.Write(val.([]byte))
	default:
		fmt.Println("default")
		err = errors.New(ERR_UNKNOW_INTERFACE_TYPE)
	}
	if err != nil {
		panic(PANIC_PUSH + err.Error())
	}
	return nb
}

func (nb *nfbuf) Pop(val interface{}) *nfbuf {
	var err error
	switch val.(type) {
	case *int:
		err = binary.Read(nb.buf, binary.LittleEndian, val)
	case *int8:
		err = binary.Read(nb.buf, binary.LittleEndian, val)
	case *int16:
		err = binary.Read(nb.buf, binary.LittleEndian, val)
	case *int32:
		err = binary.Read(nb.buf, binary.LittleEndian, val)
	case *int64:
		err = binary.Read(nb.buf, binary.LittleEndian, val)
	case *float32:
		err = binary.Read(nb.buf, binary.LittleEndian, val)
	case *float64:
		err = binary.Read(nb.buf, binary.LittleEndian, val)
	case *uint:
		err = binary.Read(nb.buf, binary.LittleEndian, val)
	case *uint8:
		err = binary.Read(nb.buf, binary.LittleEndian, val)
	case *uint16:
		err = binary.Read(nb.buf, binary.LittleEndian, val)
	case *uint32:
		err = binary.Read(nb.buf, binary.LittleEndian, val)
	case *uint64:
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

func (nb *nfbuf) Len() int {
	return nb.buf.Len()
}
