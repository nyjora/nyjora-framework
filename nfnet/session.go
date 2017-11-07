package nfnet

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"nyjora-framework/nfcommon"
)

const (
	PROTO_HEADER_SIZE = 20
	WRITE_CACHE_SIZE  = 8
)

type Protocol struct {
	Id       uint32
	FromType uint32
	FromId   uint32
	ToType   uint32
	ToId     uint32
	Data     []byte
}

type NetSession struct {
	conn       net.Conn
	reader     *bufio.Reader
	writer     *bufio.Writer
	Id         nfcommon.SessionId
	writeChan  chan []byte
	compressed bool
}

func NewNetSession(conn net.Conn) *NetSession {
	s := NetSession{
		Id:         nfcommon.NextSessionId(),
		compressed: false,
		writeChan:  make(chan []byte),
	}
	s.Reset(conn)
	return &s
}

func (ns *NetSession) String() string {
	return fmt.Sprintf("NetSession<%d@%s>", ns.Id, ns.RemoteAddr())
}

func (ns *NetSession) RemoteAddr() net.Addr {
	return ns.conn.RemoteAddr()
}

func (ns *NetSession) Run() {
	// write coroutine
	go ns.writeStream()
	// read loop
	ns.readStream()
}

func (ns *NetSession) readStream() {
	defer func() {
		ns.Close()
	}()
	for {
		// read length first
		header := make([]byte, 4)
		if _, err := io.ReadFull(ns.reader, header); err != nil {
			fmt.Printf("[NetSession] readStream: Can not read header! err = %s\n", err.Error())
			break
		}
		length := binary.LittleEndian.Uint32(header)
		if length <= PROTO_HEADER_SIZE {
			fmt.Printf("[NetSession] readStream: len is too small! len = %d\n", length)
			break
		}
		// read package
		rawData := make([]byte, length)
		if _, err := io.ReadFull(ns.reader, rawData); err != nil {
			fmt.Printf("[NetSession] readStream: Can not read enough data err = %s\n", err.Error())
			break
		}
		// build nfbuf
		pkg := nfcommon.NewNFBufBytes(rawData)
		// decode
		proto, err := ns.decode(length, pkg)
		if err != nil {
			fmt.Printf("[NetSession] readStream: Protocol decode failed !err = %s\n", err.Error())
			break
		}
		// dispatch
		ns.dispatch(proto)
	}
}

func (ns *NetSession) decode(len uint32, pkg *nfcommon.Nfbuf) (*Protocol, error) {
	//uncompress
	if ns.compressed {
		err := pkg.UnCompress()
		if err != nil {
			return nil, err
		}
	}
	// new protocol, use pool? TODO:
	proto := &Protocol{}
	// unmarshal data
	ns.UnpackProto(pkg, proto)
	return proto, nil
}

func (ns *NetSession) dispatch(proto *Protocol) {
	// notify rpc module
	ns.testEcho(proto)
}

func (ns *NetSession) writeStream() {
	for {
		rawData := <-ns.writeChan
		// 发送
		if rawData == nil {
			continue
		}
		left := len(rawData)
		fmt.Printf("[NetSession] writeStream left = %d\n", left)
		for left > 0 {
			n, err := ns.conn.Write(rawData)
			if n == left && err == nil {
				break
			}

			if n > 0 {
				rawData = rawData[n:]
				left -= n
			}

			if err != nil {
				fmt.Printf("[NetSession] writeStream err = %s\n", err.Error())
				break
			}
		}
	}
}

func (ns *NetSession) Close() {
	if err := recover(); err != nil && !IsNetError(err.(error)) {
		fmt.Printf("[NetSession] %s error : %s\n", ns.String(), err.(error))
	} else {
		fmt.Printf("[NetSession] %s disconnected.\n", ns.String())
	}
}

func (ns *NetSession) Reset(conn net.Conn) {
	ns.conn = conn
	ns.reader = bufio.NewReader(conn)
	ns.writer = bufio.NewWriter(conn)
}

func (ns *NetSession) testEcho(proto *Protocol) {
	fmt.Printf("[TestEcho] %s: %d, %d, %d, %d, %d, %s\n", ns.RemoteAddr(), proto.Id, proto.FromType, proto.FromId, proto.ToType, proto.ToId, proto.Data)
}

func (ns *NetSession) UnpackProto(nb *nfcommon.Nfbuf, proto *Protocol) {
	nb.Pop(&proto.Id).Pop(&proto.FromType).Pop(&proto.FromId).Pop(&proto.ToType).Pop(&proto.ToId)
	proto.Data = make([]byte, nb.Len())
	nb.Pop(proto.Data)
}

func (ns *NetSession) PackProto(nb *nfcommon.Nfbuf, proto *Protocol) {
	nb.Push(proto.Id).Push(proto.FromType).Push(proto.FromId).Push(proto.ToType).Push(proto.ToId)
	nb.Push(proto.Data)
}

func (ns *NetSession) Send(proto *Protocol) {
	// encode
	pkg := nfcommon.NewNFBuf()
	ns.PackProto(pkg, proto)

	// compress
	if ns.compressed {
		err := pkg.Compress()
		if err != nil {
			fmt.Printf("[NetSession] Send err = %s\n", err.Error())
			return
		}
	}
	// add len
	rawData := nfcommon.NewNFBuf()
	rawData.Push(pkg.Len()).Push(pkg.Bytes())
	// send to write loop
	ns.writeChan <- rawData.Bytes()
}
