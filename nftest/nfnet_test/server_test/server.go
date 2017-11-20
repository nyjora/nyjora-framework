package main

import (
	"fmt"
	"nyjora-framework/nfcommon"
	"nyjora-framework/nfconf"
	"nyjora-framework/nfnet"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// 用于测试的server
type GdeliveryService struct {
	server *nfnet.TcpServer
}

func NewGdeliveryServer(opt nfnet.ServerOption) *GdeliveryService {
	gs := &GdeliveryService{}
	gs.server = nfnet.NewTcpServer(opt, gs)
	return gs
}

func (gds *GdeliveryService) OnAddSession(s *nfnet.NetSession) {
	fmt.Printf("[OnAddSession] Session = %d\n", s.Id)
}

func (gds *GdeliveryService) OnDelSession(id nfcommon.SessionId) {
	fmt.Printf("[OnDelSession] Session = %d\n", id)
}

var wg *sync.WaitGroup

func ExitFunc() {
	wg.Wait()
	os.Exit(0)
}

func main() {
	fmt.Println("server_test begin...")
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

	serverOpt := nfnet.ServerOption{
		Ip:   nfconf.Conf.Get("serverconfig").Get("ip").MustString(""),
		Port: nfconf.Conf.Get("serverconfig").Get("port").MustInt(0),
	}
	wg = &sync.WaitGroup{}
	gdeliveryService := NewGdeliveryServer(serverOpt)
	go gdeliveryService.server.Serve(wg)

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("[server.go] Exit.")
				gdeliveryService.server.Stop()
				ExitFunc()
			default:
				fmt.Println("[server.go] default signal.")
			}
		}
	}()

	for {
		/*
			time.Sleep(time.Second * 10)
			fmt.Println("[server.go] Stop.")
			gdeliveryService.server.Stop()
			time.Sleep(time.Second * 10)
			fmt.Println("[server.go] Serve.")
			go gdeliveryService.server.Serve(wg)
		*/
	}
}
