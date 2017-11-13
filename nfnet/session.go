package nfnet

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"nyjora-framework/nfcommon"
	"sync"

	"github.com/golang/snappy"
)

const (
	PROTO_HEADER_SIZE       = 20
	WRITE_CACHE_SIZE        = 8
	PAYLOAD_LEN_MASK        = 0x0FFFFFFF
	PAYLOAD_COMPRESSED_MASK = 0x10000000
)

type ReadHandler interface {
	HandleProtocol(s *NetSession, p *nfcommon.Protocol)
}

type NetSession struct {
	conn        net.Conn
	reader      *bufio.Reader
	writer      *bufio.Writer
	Id          nfcommon.SessionId
	writeChan   chan []byte
	readHandler ReadHandler
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewNetSession(conn net.Conn) *NetSession {
	s := NetSession{
		Id:        nfcommon.NextSessionId(),
		writeChan: make(chan []byte, WRITE_CACHE_SIZE),
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

func (ns *NetSession) Run(wg *sync.WaitGroup) {
	// write coroutine
	go ns.writeStream(wg)
	// read loop
	ns.readStream(wg)
}

func (ns *NetSession) readStream(wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("[NetSession] readStream defer func.")
		ns.Close()
		wg.Done()
	}()
	wg.Add(1)
	for {
		// read length first
		header := make([]byte, 4)
		if _, err := io.ReadFull(ns.reader, header); err != nil {
			fmt.Printf("[NetSession] readStream: Can not read header! err = %s\n", err.Error())
			break
		}
		headerInt := binary.LittleEndian.Uint32(header)
		length := uint32(headerInt & PAYLOAD_LEN_MASK)
		if length <= PROTO_HEADER_SIZE {
			fmt.Printf("[NetSession] readStream: len is too small. len = %d\n", length)
			break
		}
		// read package
		rawData := make([]byte, length)
		if _, err := io.ReadFull(ns.reader, rawData); err != nil {
			fmt.Printf("[NetSession] readStream: Can not read enough data. err = %s\n", err.Error())
			break
		}
		// build nfbuf
		pkg := nfcommon.NewNFBufBytes(rawData)
		// decode
		compressed := (uint32(headerInt&PAYLOAD_COMPRESSED_MASK) != 0)
		proto, err := ns.decode(length, pkg, compressed)
		if err != nil {
			fmt.Printf("[NetSession] readStream: Protocol decode failed. err = %s\n", err.Error())
			break
		}
		// dispatch
		if ns.readHandler != nil {
			ns.readHandler.HandleProtocol(ns, proto)
		}

	}
}

func (ns *NetSession) writeStream(wg *sync.WaitGroup) {
	fmt.Println("[NetSession] writeStream begin.")
	wg.Add(1)
	defer func() {
		fmt.Println("[NetSession] writeStream defer func.")
		wg.Done()
	}()
	for {
		select {
		case <-ns.ctx.Done():
			fmt.Println("[NetSession] writeStream ctx.Done().")
			return
		case rawData := <-ns.writeChan:
			// 发送
			if rawData == nil {
				continue
			}
			left := len(rawData)
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
}

func (ns *NetSession) decode(length uint32, pkg *nfcommon.Nfbuf, compressed bool) (*nfcommon.Protocol, error) {
	//uncompress
	var cbuf []byte
	var err error
	if compressed {
		cbuf, err = snappy.Decode(nil, pkg.Bytes())
		if err != nil {
			return nil, err
		}
	} else {
		cbuf = pkg.Bytes()
	}
	apkg := nfcommon.NewNFBufBytes(cbuf)
	// new protocol, use pool? TODO:
	proto := nfcommon.NewProto()
	// unmarshal data
	ns.UnpackProto(apkg, proto)
	return proto, nil
}

func (ns *NetSession) encode(proto *nfcommon.Protocol) *nfcommon.Nfbuf {
	pkg := nfcommon.NewNFBuf()
	ns.PackProto(pkg, proto)
	cbuf := snappy.Encode(nil, pkg.Bytes())
	rawData := nfcommon.NewNFBuf()
	if cbuf != nil && int32(len(cbuf)) < pkg.Len() {
		pkg = nfcommon.NewNFBufBytes(cbuf)
		length := pkg.Len()
		length = length&0x0FFFFFFF | PAYLOAD_COMPRESSED_MASK
		rawData.Push(length).Push(pkg.Bytes())
		return rawData
	}
	rawData.Push(pkg.Len() & 0x0FFFFFFF).Push(pkg.Bytes())
	return rawData
}

func (ns *NetSession) Close() {
	fmt.Println("[NetSession] Close.")
	if err := recover(); err != nil && !IsNetError(err.(error)) {
		fmt.Printf("[NetSession] %s error : %s\n", ns.String(), err.(error))
	} else {
		fmt.Printf("[NetSession] %s disconnected.\n", ns.String())
	}
	ns.conn.Close()
	ns.cancel()
}

func (ns *NetSession) Reset(conn net.Conn) {
	ns.conn = conn
	ns.reader = bufio.NewReader(conn)
	ns.writer = bufio.NewWriter(conn)
	ns.ctx, ns.cancel = context.WithCancel(context.Background())
}

func (ns *NetSession) UnpackProto(nb *nfcommon.Nfbuf, proto *nfcommon.Protocol) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	nb.Pop(&proto.Id).Pop(&proto.FromType).Pop(&proto.FromId).Pop(&proto.ToType).Pop(&proto.ToId)
	proto.Data = make([]byte, nb.Len())
	nb.Pop(proto.Data)
}

func (ns *NetSession) PackProto(nb *nfcommon.Nfbuf, proto *nfcommon.Protocol) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	nb.Push(proto.Id).Push(proto.FromType).Push(proto.FromId).Push(proto.ToType).Push(proto.ToId)
	nb.Push(proto.Data)
}

func (ns *NetSession) Send(proto *nfcommon.Protocol) {
	// encode
	rawData := ns.encode(proto)
	// send to write loop
	select {
	case ns.writeChan <- rawData.Bytes():
	default:
		fmt.Println("[NetSession] Send writeChan cache full!")
		return
	}
}

func (ns *NetSession) RegisterNub(rh ReadHandler) {
	ns.readHandler = rh
}
