package main

import (
	"fmt"
	"nyjora-framework/nfcommon"
	"nyjora-framework/nfconf"
	"nyjora-framework/nfnet"
	"os"
	"time"
)

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

func main() {
	fmt.Println("server_test begin...")
	filepath := os.Args[1]
	fmt.Println("Read Json file : " + filepath)
	err := nfconf.Init(filepath)
	if err != nil {
		fmt.Printf("[Main] conf init err! : %s\n", err)
		os.Exit(1)
	}
	nfconf.Conf.Get("testconfig").Get("port").Int()
	serverOpt := nfnet.ServerOption{
		Ip:   nfconf.Conf.Get("serverconfig").Get("ip").MustString(""),
		Port: nfconf.Conf.Get("serverconfig").Get("port").MustInt(0),
	}

	tserver := &TcpServerDelegate{}
	tserver.Init(serverOpt)
	go tserver.Serve()
	for {
		time.Sleep(10000)
	}
}
