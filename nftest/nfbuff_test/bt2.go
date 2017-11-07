package main

import (
	"fmt"
	"nyjora-framework/nfcommon"
)

func main() {
	buf := nfcommon.NewNFBuf()
	var t1 int = 1
	var t2 int32 = 2
	var t3 int8 = 3
	var t4 float32 = 4.1
	var t5 float64 = 5.2
	t6 := []byte{'f', 'u', 'c', 'k'}
	var t7 byte = 'c'
	buf.Push(t1).Push(t2).Push(t3).Push(t4).Push(t5).Push(t6).Push(t7)
	//buf.Push(t2)
	//buf.Push(t3)
	//buf.Push(t4)
	//buf.Push(t5)
	//buf.Push(t6)
	//buf.Push(t7)

	fmt.Printf("buf size = %d\n", buf.Len())

	var tt1 int
	var tt2 int32
	var tt3 int8
	var tt4 float32
	var tt5 float64
	tt6 := make([]byte, 4)
	var tt7 byte

	buf.Pop(&tt1).Pop(&tt2).Pop(&tt3).Pop(&tt4).Pop(&tt5).Pop(tt6).Pop(&tt7)
	//buf.Pop(&tt2)
	//buf.Pop(&tt3)
	//buf.Pop(&tt4)
	//buf.Pop(&tt5)
	//buf.Pop(tt6)
	//buf.Pop(&tt7)

	fmt.Printf("tt1 = %d\n", tt1)
	fmt.Printf("tt2 = %d\n", tt2)
	fmt.Printf("tt3 = %d\n", tt3)
	fmt.Printf("tt4 = %g\n", tt4)
	fmt.Printf("tt5 = %g\n", tt5)
	fmt.Printf("tt6 = %s\n", tt6)
	fmt.Printf("tty = %c\n", tt7)
}
