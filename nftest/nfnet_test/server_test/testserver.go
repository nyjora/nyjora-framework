package server_test

type TcpServerDelegate struct {
	tserver *nfnet.TcpServer
}

func (tsd *TcpServerDelegate) Init(opts nfnet.ServerOption) {
	tsd.tserver = nfnet.NewTcpServer(opts, tsd)
}

func (tsd *TcpServerDelegate) OnAddSession(id nfcommon.SessionId) {
	fmt.Printf("[OnAddSession] Session = %d\n", id)
}

func (tsd *TcpServerDelegate) OnDelSession(id nfcommon.SessionId) {
	fmt.Printf("[OnDelSession] Session = %d\n", id)
}

func (tsd *TcpServerDelegate) Serve() {
	if tsd.tserver != nil {
		tsd.tserver.Run()
	}
}

var service TcpServerDelegate

