package nfnet

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"nyjora-framework/nfcommon"

	"github.com/golang/snappy"
)

const (
	MSG_HEADER_SIZE = 20
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
	compressed bool
}

func NewNetSession(conn net.Conn) *NetSession {
	s := NetSession{
		Id:         nfcommon.NextSessionId(),
		compressed: false,
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
	// read loop
	ns.readStream()
	// write coroutine
	go ns.writeStream()
}

func (ns *NetSession) readStream() {
	defer func() {
		ns.Close()
	}()
	for {
		// read length first
		header := make([]byte, 2)
		if _, err := io.ReadFull(ns.reader, header); err != nil {
			fmt.Printf("[NetSession] readStream: Can not read header! err = %s\n", err.Error())
			break
		}
		len := binary.LittleEndian.Uint16(header)
		if len <= MSG_HEADER_SIZE {
			fmt.Printf("[NetSession] readStream: len is too small! len = %d\n", len)
			break
		}
		// read package
		pkg := make([]byte, len)
		if _, err := io.ReadFull(ns.reader, pkg); err != nil {
			fmt.Printf("[NetSession] readStream: Can not read enough data err = %s\n", err.Error())
			break
		}
		// decode
		proto, derr := ns.decode(len, pkg)
		if derr != nil {
			fmt.Printf("[NetSession] readStream: Protocol decode failed !err = %s\n", derr.Error())
			break
		}
		// dispatch
		ns.dispatch(proto)
	}
}

func (ns *NetSession) decode(len uint16, pkg []byte) (*Protocol, error) {
	//uncompress
	var dpkg []byte
	var err error
	if ns.compressed {
		dpkg, err = snappy.Decode(nil, pkg)

		if err != nil {
			return nil, err
		}
	} else {
		dpkg = pkg
	}

	// new protocol, use pool? TODO:
	proto := Protocol{}
	// unmarshal data
	proto.Id = binary.LittleEndian.Uint32(dpkg[:4])
	proto.FromType = binary.LittleEndian.Uint32(dpkg[4:8])
	proto.FromId = binary.LittleEndian.Uint32(dpkg[8:12])
	proto.ToType = binary.LittleEndian.Uint32(dpkg[12:16])
	proto.ToId = binary.LittleEndian.Uint32(dpkg[16:20])
	proto.Data = dpkg[20:]

	return &proto, nil
}

func (ns *NetSession) dispatch(proto *Protocol) {
	// notify rpc module
	ns.testEcho(proto)
}

func (ns *NetSession) writeStream() {

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
	fmt.Printf("[TestEcho] %s: %d, %d, %d, %d, %d, %s\n", ns.RemoteAddr(), proto.Id, proto.FromType, proto.FromId, proto.ToType, proto.ToId)
}
