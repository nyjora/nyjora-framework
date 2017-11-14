package nfdb

import (
	"context"
	"fmt"
	"nyjora-framework/nfconf"
	"sync"
	"time"

	"github.com/robfig/cron"

	"github.com/boltdb/bolt"
)

var DB *bolt.DB
var Ctx context.Context
var cancel context.CancelFunc
var Wg *sync.WaitGroup
var hbCron *cron.Cron

func Start() error {
	// Initialze database config, MUST after conf init !!
	dbname := nfconf.Conf.Get("database").Get("name").MustString("bolt.db")
	dbpath := nfconf.Conf.Get("database").Get("path").MustString("./")
	dbtables := nfconf.Conf.Get("database").Get("tables").MustStringArray()
	var err error
	DB, err = bolt.Open(dbpath+dbname, 0600, nil)
	if err != nil {
		fmt.Printf("[boltdb.go] db open fail. name = %s\n", dbpath+dbname)
		return err
	}
	fmt.Printf("[boltdb.go] Start. %s\n", dbpath+dbname)
	// Initialse database tables(buckets)
	for _, v := range dbtables {
		err = DB.Update(func(tx *bolt.Tx) error {
			fmt.Printf("[boltdb.go] CreateBucketIfNotExists v = %v\n", v)
			_, err := tx.CreateBucketIfNotExists([]byte(v))
			if err != nil {
				fmt.Printf("[boltdb.go] create table %s fail.\n", v)
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	// begin hotbackup cron
	hbCron = cron.New()
	hbCron.AddFunc("0 30 3 * * *", hotBackup)
	hbCron.Start()

	// new wg
	Wg = &sync.WaitGroup{}
	// context
	Ctx, cancel = context.WithCancel(context.Background())
	return nil
}

func Close() {
	fmt.Println("[boltdb.go] Close.")
	hbCron.Stop()
	cancel()
	Wg.Wait()
	if DB != nil {
		DB.Close()
		DB = nil
	}
}

func hotBackup() {
	backup := nfconf.Conf.Get("database").Get("backup").MustString("./")
	timestring := time.Now().Format("20060102150405")
	if DB != nil {
		fmt.Printf("[boltdb.go] hotbackup. %s\n", backup)
		err := DB.View(func(tx *bolt.Tx) error {
			return tx.CopyFile(backup+timestring+".dbbk", 0666)
		})
		if err != nil {
			fmt.Printf("[boltdb.go] hotBackup failed at %s\n")
			return
		}
	}
}

func BackUp() {
	hotBackup()
}
