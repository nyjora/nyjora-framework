package main

import (
	"errors"
	"fmt"
	"nyjora-framework/nfconf"
	"nyjora-framework/nfdb"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boltdb/bolt"
)

func TestDBWrite() {
	fmt.Println("[TestDBWrite]---")
	nfdb.Wg.Add(1)
	defer nfdb.Wg.Done()
	// add data to base
	select {
	case <-nfdb.Ctx.Done():
		fmt.Println("[TestDBWrite]")
		return
	default:
		err := nfdb.DB.Update(func(tx *bolt.Tx) error {
			err := tx.Bucket([]byte("base")).Put([]byte("chenyan"), []byte("mage"))
			err = tx.Bucket([]byte("base")).Put([]byte("shideqiang"), []byte("ranger"))
			err = tx.Bucket([]byte("status")).Put([]byte("chenyan"), []byte("100"))
			return err
		})
		if err != nil {
			fmt.Printf("[TestDBWrite] write error %v\n", err)
		}
		return
	}
}

func TestDBRead() {
	fmt.Println("[TestDBRead]---")
	nfdb.Wg.Add(1)
	defer nfdb.Wg.Done()
	// add data to base
	select {
	case <-nfdb.Ctx.Done():
		fmt.Println("[TestDBRead]")
		return
	default:
		err := nfdb.DB.View(func(tx *bolt.Tx) error {
			val := tx.Bucket([]byte("base")).Get([]byte("chenyan"))
			if val == nil {
				return errors.New("Db not found.")
			}
			fmt.Printf("[base]chenyan : %s\n", val)
			val2 := tx.Bucket([]byte("status")).Get([]byte("chenyan"))
			if val2 == nil {
				return errors.New("Db not found.")
			}
			fmt.Printf("[status]chenyan : %s\n", val2)
			return nil
		})
		if err != nil {
			fmt.Printf("[TestDBRead] read error %v\n", err)
		}
		return
	}
}

func TestDBDel() {
	fmt.Println("[TestDBDel]---")
	nfdb.Wg.Add(1)
	defer nfdb.Wg.Done()
	// add data to base
	select {
	case <-nfdb.Ctx.Done():
		fmt.Println("[TestDBDel]")
		return
	default:
		err := nfdb.DB.Update(func(tx *bolt.Tx) error {
			err := tx.Bucket([]byte("base")).Delete([]byte("chenyan"))
			return err
		})
		if err != nil {
			fmt.Printf("[TestDBDel] delete error %v\n", err)
		}
		return
	}
}

func ExitFunc() {
	fmt.Println("[ExitFunc]")
	os.Exit(0)
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

	nfdb.Start()
	//go TestDBWrite()
	go TestDBRead()
	/*
		time.Sleep(100 * time.Microsecond)
		go TestDBRead()
		time.Sleep(100 * time.Microsecond)
		go TestDBDel()
		time.Sleep(100 * time.Microsecond)
		go TestDBRead()
		time.Sleep(100 * time.Microsecond)
		go TestDBWrite()
		time.Sleep(100 * time.Microsecond)
		go TestDBRead()
		// DB test
	*/

	time.Sleep(5 * time.Second)
	nfdb.BackUp()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("[client.go] Exit.")
				nfdb.Close()
				ExitFunc()
			default:
				fmt.Println("[client.go] default signal.")
			}
		}
	}()

	for {
	}
}
