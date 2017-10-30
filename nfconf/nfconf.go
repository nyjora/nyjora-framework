package nfconf

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// 映射配置文件的结构体，每个结构体针对一个文件
type NetConfig struct {
	TestConfig struct {
		Ip   string `json:"ip"`
		Port int    `json:"port"`
	} `json:"testconfig`
	Name string `json:"name"`
}

var NetConf NetConfig
var once sync.Once

func Init(filePath string) (e error) {
	var initFunc = func() {
		confFile, err := os.Open(filePath)
		defer confFile.Close()
		if err != nil {
			fmt.Println(err.Error())
			e = err
			return
		}
		err = json.NewDecoder(confFile).Decode(&NetConf)
		if err != nil {
			fmt.Printf("Init json config file failed,%s,%s\n", filePath, err)
			e = err
			return
		}
		e = err
	}
	once.Do(initFunc)
	return e
}
