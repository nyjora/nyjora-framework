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

func (tcd *GDeliveryClient) Connect(wg *sync.WaitGroup) {
	if tcd.client != nil {
		tcd.client.Run(wg)
	}
}

func ExitFunc() {
	wg.Wait()
	os.Exit(0)
}

var wg *sync.WaitGroup

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
	wg = &sync.WaitGroup{}
	gdeliveryClient := NewGDeliveryClient(clientOpt)
	gdeliveryClient.Connect(wg)

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("[server.go] Exit.")
				gdeliveryClient.client.Stop(wg)
				ExitFunc()
			default:
				fmt.Println("[server.go] default signal.")
			}
		}
	}()

	for {
		time.Sleep(time.Second * 5)
	}
}
