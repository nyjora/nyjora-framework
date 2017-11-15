package main

import (
	"nyjora-framework/example/gdatabase/service"
	"nyjora-framework/nfconf"
	"nyjora-framework/nfdb"
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
	nflog.Init(os.Stdout, os.Stdout, os.Stdout, os.Stdout, "gdatabase")
}

func main() {
	InitLogger()
	InitConf()
	wg = &sync.WaitGroup{}

	// start db
	nfdb.Start()

	// start dbserver
	gdatabase.InitDBServer(nfnet.ServerOption{
		Ip:   nfconf.Conf.Get("dbserver").Get("ip").MustString(""),
		Port: nfconf.Conf.Get("dbserver").Get("port").MustInt(0),
	})
	go gdatabase.DBServer.Run(wg)

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				nflog.Debug("Exit...")
				gdatabase.DBServer.Stop(wg)
				nfdb.Close()
				Exit()
			default:
				nflog.Err("unknown signal.")
			}
		}
	}()

	for {
		// TODO:
	}
}
