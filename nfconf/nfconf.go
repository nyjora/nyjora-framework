package nfconf

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type NetConfig struct {
	TestConfig struct {
		Ip   string `json:"ip"`
		Port int    `json:"port"`
	} `json:"testconfig`
	Name string `json:"name"`
}

var NetConf *NetConfig
var once sync.Once

func Init(filePath string) {
	var initFunc = func() {
		confFile, err := os.Open(filePath)
		defer confFile.Close()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		var configData NetConfig
		error := json.NewDecoder(confFile).Decode(&configData)
		if error != nil {
			fmt.Printf("Init json config file failed,%s,%s\n", filePath, error)
		}
		NetConf = &configData
	}
	once.Do(initFunc)
}
