package globby

import (
	"fmt"
	"nyjora-framework/nfcommon"
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
	fmt.Printf("[OnAddSession] Session = %d\n", id)
}

func (bc *BusClientDelegate) OnDelSession(id nfcommon.SessionId) {
	fmt.Printf("[OnDelSession] Session = %d\n", id)
}
