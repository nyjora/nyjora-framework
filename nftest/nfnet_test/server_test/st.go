package main

import (
	"fmt"
	"nyjora-framework/nfcommon"
	"nyjora-framework/nfconf"
	"nyjora-framework/nfnet"
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

	/*
	tserver := &TcpServerDelegate{}
	tserver.Init(serverOpt)
	go tserver.Serve()
	for {
		time.Sleep(10000)
	}
	*/
}