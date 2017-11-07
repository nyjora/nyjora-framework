package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

func check(x interface{}) {
	switch x.(type) {
	case int:
		fmt.Println("int")
	case int16:
		fmt.Println("int16")
	case int32:
		fmt.Println("int32")
	case string:
		fmt.Println("string")
	case []byte:
		fmt.Println("[]byte")
	case *int16:
		fmt.Println("pointer")
	default:
		fmt.Println("default")
	}
}
func main() {
	buf := new(bytes.Buffer)
	var t1 uint8 = 1
	var t2 int8 = 2
	var t3 int64 = 3
	var t4 float32 = 4.1
	var t5 float64 = 5.2
	t6 := []byte{'f', 'u', 'c', 'k'}
	var t7 byte = 'c'
	binary.Write(buf, binary.LittleEndian, t1)
	binary.Write(buf, binary.LittleEndian, t2)
	binary.Write(buf, binary.LittleEndian, t3)
	binary.Write(buf, binary.LittleEndian, t4)
	binary.Write(buf, binary.LittleEndian, t5)
	buf.Write(t6)
	buf.WriteByte(t7)
	fmt.Printf("buf size = %d\n", buf.Len())
	check(t1)
	fmt.Printf("t6 type = %s\n", reflect.TypeOf(t6))
	var tt1 int8
	var tt2 uint8
	var tt3 int64
	var tt4 float32
	var tt5 float64
	tt6 := make([]byte, 4)
	var tt7 byte
	binary.Read(buf, binary.LittleEndian, &tt1)
	binary.Read(buf, binary.LittleEndian, &tt2)
	binary.Read(buf, binary.LittleEndian, &tt3)
	binary.Read(buf, binary.LittleEndian, &tt4)
	binary.Read(buf, binary.LittleEndian, &tt5)
	buf.Read(tt6)
	tt7, _ = buf.ReadByte()

	fmt.Printf("tt1 = %d\n", tt1)
	fmt.Printf("tt2 = %d\n", tt2)
	fmt.Printf("tt3 = %d\n", tt3)
	fmt.Printf("tt4 = %g\n", tt4)
	fmt.Printf("tt5 = %g\n", tt5)
	fmt.Printf("tt6 = %s\n", tt6)
	fmt.Printf("tty = %c\n", tt7)

}
