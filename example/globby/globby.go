package main

import (
	"fmt"
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
}

func InitLogger() {
	nflog.Init(os.Stdout, os.Stdout, os.Stdout, os.Stdout, "globby")
}

func main() {
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
				fmt.Println("[globby] Exit.")
				globby.BusClient.Stop(wg)
				Exit()
			default:
				fmt.Println("[globby] default signal.")
			}
		}
	}()

	for {
		// TODO:
	}
}
