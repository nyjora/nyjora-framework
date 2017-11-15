package gdatabase

import (
	"fmt"
	"nyjora-framework/nfcommon"
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
	fmt.Printf("[OnAddSession] Session = %d\n", s.Id)
}

func (ds *DBServerDelegate) OnDelSession(id nfcommon.SessionId) {
	fmt.Printf("[OnDelSession] Session = %d\n", id)
}
