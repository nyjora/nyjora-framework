package nfnet

import (
	"fmt"
	"net"
	"nyjora-framework/nfcommon"
	"sync"
	"time"
)

const (
	RESTART_TCP_CLIENT_INTERVAL = 5 * time.Second
	STREAM_WRITE_BUFFER_SIZE    = 1048576
	STREAM_READ_BUFFER_SIZE     = 1048576
)

type ClientDelegate interface {
	OnAddSession(nfcommon.SessionId)
	OnDelSession(nfcommon.SessionId)
}

type ClientOption struct {
	Ip   string
	Port int
}

type TcpClient struct {
	opts      ClientOption
	session   *NetSession
	delegate  ClientDelegate
	connected bool
	wg        *sync.WaitGroup
}

func NewTcpClient(opt ClientOption, d ClientDelegate) *TcpClient {
	return &TcpClient{
		opts:      opt,
		delegate:  d,
		connected: false,
	}
}

func (tc *TcpClient) Run(wg *sync.WaitGroup) {
	wg.Add(1)
	for {
		if tc.connect(wg) {
			tc.wg = &sync.WaitGroup{}
			tc.session.Run(tc.wg)
		}
		time.Sleep(RESTART_TCP_CLIENT_INTERVAL)
	}
}

func (tc *TcpClient) connect(wg *sync.WaitGroup) bool {
	addr := fmt.Sprintf("%s:%d", tc.opts.Ip, tc.opts.Port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("[TcpClient] Connect to server failed. addr = %s\n", addr)
		return false
	}
	if conn == nil {
		fmt.Printf("[TcpClient] conn is nil. addr = %s\n", addr)
		return false
	}
	tc.connected = true
	if tc.session == nil {
		tc.session = NewNetSession(conn)
	} else {
		// reconnect
		tc.session.Reset(conn)
	}
	fmt.Printf("[TcpClient] Server connected! addr = %s\n", addr)
	return true
}

func (tc *TcpClient) SendProto(id uint32, fromType uint32, fromId uint32, toType uint32, toId uint32, data []byte) {
	if tc.session == nil || tc.connected == false {
		return
	}
	proto := nfcommon.NewProto()
	proto.Id = id
	proto.FromType = fromType
	proto.FromId = fromId
	proto.ToType = toType
	proto.ToId = toId
	proto.Data = data
	tc.session.Send(proto)
}

func (tc *TcpClient) IsValid() bool {
	if tc.session == nil || tc.connected == false {
		return false
	}
	return true
}

func (tc *TcpClient) Stop(wg *sync.WaitGroup) {
	fmt.Println("[TcpClient] Stop.")
	if tc.session != nil {
		tc.session.Close()
		tc.session = nil
		tc.wg.Wait()
	}
	wg.Done()
}
