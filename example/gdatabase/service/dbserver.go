package gdatabase

import (
	"nyjora-framework/nfcommon"
	"nyjora-framework/nflog"
	"nyjora-framework/nfnet"
)

type DBServerDelegate struct {
	*nfnet.TcpServer
}

var DBServer *DBServerDelegate

func InitDBServer(opt nfnet.ServerOption) {
	DBServer = &DBServerDelegate{
		nfnet.NewTcpServer(opt, DBServer),
	}
}

func (ds *DBServerDelegate) OnAddSession(s *nfnet.NetSession) {
	nflog.Info("[OnAddSession] Session = %d\n", s.Id)
}

func (ds *DBServerDelegate) OnDelSession(id nfcommon.SessionId) {
	nflog.Info("[OnDelSession] Session = %d\n", id)
}
