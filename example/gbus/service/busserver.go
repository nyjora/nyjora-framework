package gbus

import (
	"fmt"
	"nyjora-framework/nfcommon"
	"nyjora-framework/nfnet"
)

type BusServerDelegate struct {
	*nfnet.TcpServer
}

var BusServer *BusServerDelegate

func InitBusServer(opt nfnet.ServerOption) {
	BusServer = &BusServerDelegate{
		nfnet.NewTcpServer(opt, BusServer),
	}
}

func (bs *BusServerDelegate) OnAddSession(s *nfnet.NetSession) {
	fmt.Printf("[OnAddSession] Session = %d\n", s.Id)
}

func (bs *BusServerDelegate) OnDelSession(id nfcommon.SessionId) {
	fmt.Printf("[OnDelSession] Session = %d\n", id)
}
