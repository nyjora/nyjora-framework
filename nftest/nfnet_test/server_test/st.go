package main

import (
	"fmt"
	"nyjora-framework/nfconf"
	"nyjora-framework/nfnet"
	"nyjora-framework/nfproto"
	"nyjora-framework/nftest/nflogic_test"
	"nyjora-framework/nftest/nfservice_test"
	"os"
	"time"
)

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
	nfproto.Init()
	nflogic_test.Nftestserver = nfservice_test.NewTcpServerDelegate()
	nflogic_test.Nftestserver.Init(serverOpt)
	go nflogic_test.Nftestserver.Serve()
	for {
		time.Sleep(time.Second * 60)
	}

}
