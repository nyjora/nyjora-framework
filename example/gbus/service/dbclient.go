package gbus

import (
	"fmt"
	"nyjora-framework/nfcommon"
	"nyjora-framework/nfnet"
)

type DBClientDelegate struct {
	*nfnet.TcpClient
}

var DBClient *DBClientDelegate

func InitDBClient(opt nfnet.ClientOption) {
	DBClient = &DBClientDelegate{
		nfnet.NewTcpClient(opt, DBClient),
	}
}

func (dc *DBClientDelegate) OnAddSession(id nfcommon.SessionId) {
	fmt.Printf("[OnAddSession] Session = %d\n", id)
}

func (dc *DBClientDelegate) OnDelSession(id nfcommon.SessionId) {
	fmt.Printf("[OnDelSession] Session = %d\n", id)
}
