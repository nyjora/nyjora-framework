package nfnet

import (
	"fmt"
	"net"
	"nyjora-framework/nfcommon"
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
	opts     ClientOption
	session  *NetSession
	delegate ClientDelegate
}

func NewTcpClient(opt ClientOption, d ClientDelegate) *TcpClient {
	return &TcpClient{
		opts:     opt,
		delegate: d,
	}
}

func (tc *TcpClient) Run() {
	addr := fmt.Sprintf("%s:%d", tc.opts.Ip, tc.opts.Port)
	go tc.connect(addr)
}

func (tc *TcpClient) connect(addr string) {
	fmt.Printf("[TcpClient] connect coroutine begin! addr = %s\n", addr)
	defer func() {
		if tc.session != nil {
			tc.session.Close()
			tc.session = nil
		}
	}()
	fmt.Printf("[TcpClient] connect coroutine begin! addr = %s\n", addr)
	for {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Printf("[TcpClient] Connect to server failed. addr = %s\n", addr)
			return
		}
		if conn == nil {
			fmt.Printf("[TcpClient] conn is nil. addr = %s\n", addr)
			return
		}
		if tc.session == nil {
			tc.session = NewNetSession(conn)
		} else {
			tc.session.Reset(conn)
		}
		fmt.Printf("[TcpClient] Server connected! addr = %s\n", addr)
		tc.session.Run()
	}
}
