package main

import (
	"fmt"
	"nyjora-framework/nfcommon"
	"nyjora-framework/nfconf"
	"nyjora-framework/nfnet"
	"os"
	"time"
)

type TcpClientDelegate struct {
	tclient *nfnet.TcpClient
}

func (tcd *TcpClientDelegate) Init(opts nfnet.ClientOption) {
	tcd.tclient = nfnet.NewTcpClient(opts, tcd)
}

func (tcd *TcpClientDelegate) OnAddSession(id nfcommon.SessionId) {
	fmt.Printf("[OnAddSession] Session = %d\n", id)
}

func (tcd *TcpClientDelegate) OnDelSession(id nfcommon.SessionId) {
	fmt.Printf("[OnDelSession] Session = %d\n", id)
}

func (tcd *TcpClientDelegate) Connect() {
	fmt.Println("[TcpClientDelegate] Connect!")
	if tcd.tclient != nil {
		tcd.tclient.Run()
	}
}

func main() {
	fmt.Println("client begin...")
	filepath := os.Args[1]
	fmt.Println("Read Json file : " + filepath)
	err := nfconf.Init(filepath)
	if err != nil {
		fmt.Printf("[Main] conf init err! : %s\n", err)
		os.Exit(1)
	}
	clientOpt := nfnet.ClientOption{
		Ip:   nfconf.Conf.Get("clientconfig").Get("ip").MustString(),
		Port: nfconf.Conf.Get("clientconfig").Get("port").MustInt(),
	}

	tclient := &TcpClientDelegate{}
	tclient.Init(clientOpt)
	tclient.Connect()
	msg := []byte{'f', 'u', 'c', 'k'}
	for {
		time.Sleep(time.Second * 5)
		if tclient.tclient.IsValid() {
			fmt.Println("Begin to send...")
			tclient.tclient.SendProto(1, 2, 3, 4, 5, msg)
		}
	}
}
