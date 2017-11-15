package main

import (
	globby "nyjora-framework/example/globby/service"
	"nyjora-framework/nfconf"
	"nyjora-framework/nflog"
	"nyjora-framework/nfnet"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// global waitgroup
var wg *sync.WaitGroup

func Exit() {
	wg.Wait()
	os.Exit(0)
}

func InitConf() {
	if len(os.Args) < 2 {
		nflog.Fatal("Args too short %v\n", len(os.Args))
	}
	filepath := os.Args[1]
	nflog.Info("Read Json file : " + filepath)
	err := nfconf.Init(filepath)
	if err != nil {
		nflog.Fatal("conf init err : %v\n", err)
	}
}

func InitLogger() {
	nflog.Init(os.Stdout, os.Stdout, os.Stdout, os.Stdout, "globby")
}

func main() {
	InitLogger()
	InitConf()
	wg = &sync.WaitGroup{}

	// start busclient
	globby.InitDBClient(nfnet.ClientOption{
		Ip:   nfconf.Conf.Get("busserver").Get("ip").MustString(""),
		Port: nfconf.Conf.Get("busserver").Get("port").MustInt(0),
	})
	go globby.BusClient.Run(wg)

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				nflog.Debug("Exit...")
				globby.BusClient.Stop(wg)
				Exit()
			default:
				nflog.Err("unknow signal.")
			}
		}
	}()

	for {
		// TODO:
	}
}
