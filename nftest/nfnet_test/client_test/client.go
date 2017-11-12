package main

import (
	"fmt"
	"nyjora-framework/nfcommon"
	"nyjora-framework/nfconf"
	"nyjora-framework/nfnet"
	"os"
	"time"
)

type GDeliveryClient struct {
	client *nfnet.TcpClient
}

func NewGDeliveryClient(opt nfnet.ClientOption) *GDeliveryClient {
	gc := &GDeliveryClient{}
	gc.client = nfnet.NewTcpClient(opt, gc)
	return gc
}

func (tcd *GDeliveryClient) OnAddSession(id nfcommon.SessionId) {
	fmt.Printf("[OnAddSession] Session = %d\n", id)
}

func (tcd *GDeliveryClient) OnDelSession(id nfcommon.SessionId) {
	fmt.Printf("[OnDelSession] Session = %d\n", id)
}

func (tcd *GDeliveryClient) Connect() {
	if tcd.client != nil {
		tcd.client.Run()
	}
}

func main() {
	fmt.Println("client begin...")
	if len(os.Args) < 2 {
		fmt.Println("Args too short!")
		os.Exit(1)
	}
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
	gdeliveryClient := NewGDeliveryClient(clientOpt)
	gdeliveryClient.Connect()
	for {
		time.Sleep(time.Second * 5)
	}
}
