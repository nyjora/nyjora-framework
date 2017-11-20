package nfservice_test

import (
	"fmt"
	"nyjora-framework/nfcommon"
	"nyjora-framework/nfnet"
	"nyjora-framework/nfproto"
)

type TcpServerDelegate struct {
	tserver *nfnet.TcpServer
}

func NewTcpServerDelegate() *TcpServerDelegate {
	return &TcpServerDelegate{}
}

func (tsd *TcpServerDelegate) Init(opts nfnet.ServerOption) {
	tsd.tserver = nfnet.NewTcpServer(opts, tsd)
}

func (tsd *TcpServerDelegate) OnAddSession(s *nfnet.NetSession) {
	nub := nfproto.NewDispatchTest()
	s.RegisterNub(nub)
	fmt.Printf("[OnAddSession] Session = %d\n", s.Id)
}

func (tsd *TcpServerDelegate) OnDelSession(id nfcommon.SessionId) {
	fmt.Printf("[OnDelSession] Session = %d\n", id)
}

func (tsd *TcpServerDelegate) Serve() {
	if tsd.tserver != nil {
		tsd.tserver.Run()
	}
}

func (tsd *TcpServerDelegate) BroadCastChatTest(id uint32, fromType uint32, fromId uint32, toType uint32, toId uint32, data []byte) {
	tsd.tserver.BroadCast(id, fromType, fromId, toType, toId, data)
}
