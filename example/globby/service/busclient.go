package globby

import (
	"nyjora-framework/nfcommon"
	"nyjora-framework/nflog"
	"nyjora-framework/nfnet"
)

type BusClientDelegate struct {
	*nfnet.TcpClient
}

var BusClient *BusClientDelegate

func InitDBClient(opt nfnet.ClientOption) {
	BusClient = &BusClientDelegate{
		nfnet.NewTcpClient(opt, BusClient),
	}
}

func (bc *BusClientDelegate) OnAddSession(id nfcommon.SessionId) {
	nflog.Info("[OnAddSession] Session = %d\n", id)
}

func (bc *BusClientDelegate) OnDelSession(id nfcommon.SessionId) {
	nflog.Info("[OnDelSession] Session = %d\n", id)
}
