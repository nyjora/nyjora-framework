package gbus

import (
	"nyjora-framework/nfcommon"
	"nyjora-framework/nflog"
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
	nflog.Info("[OnAddSession] Session = %d\n", s.Id)
}

func (bs *BusServerDelegate) OnDelSession(id nfcommon.SessionId) {
	nflog.Info("[OnDelSession] Session = %d\n", id)
}
