package main

import (
	"fmt"
	"nyjora-framework/nfconf"
	"os"
)

func main() {
	filepath := os.Args[1]
	fmt.Println("Read Json file : " + filepath)
	err := nfconf.Init(filepath)
	if err != nil {
		fmt.Printf("[Main] Catch err! err  : %s\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("[Main] Init no err!\n")
	}
	port, _ := nfconf.Conf.Get("testconfig").Get("port").Int()
	ip, _ := nfconf.Conf.Get("testconfig").Get("ip").String()
	name, _ := nfconf.Conf.Get("name").String()
	fmt.Printf("Config port = %v\n", port)
	fmt.Printf("Config ip = %v\n", ip)
	fmt.Printf("Config name = %v\n", name)

	fmt.Printf("Must port = %v\n", nfconf.Conf.Get("testconfig").Get("port").MustInt())
	fmt.Printf("Must ip = %v\n", nfconf.Conf.Get("testconfig").Get("ip").MustString("no ip"))
	fmt.Printf("Must name = %v\n", nfconf.Conf.Get("name").MustString("no name"))
}
