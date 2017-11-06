package nfcommon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type nfbuf struct {
	buf *bytes.Buffer
}

func NewNFBuf() *nfbuf {
	return &nfbuf{
		buf: new(bytes.Buffer),
	}
}

func (nb *nfbuf) Push(val interface{}) (err error) {
	fmt.Println("push")
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
		err = errors.New("[nfbuf] Push unknow type.")
	}
	return err
}

func (nb *nfbuf) Pop(val interface{}) (err error) {
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
	return err
}

func (nb *nfbuf) Len() int {
	return nb.buf.Len()
}
