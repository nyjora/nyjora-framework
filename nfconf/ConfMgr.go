package nfconf

import "os"
import "fmt"
import "encoding/json"

// 所有的config都保存在此，并由ConfMgr统一管理起来
type TestConfData struct {
	Ip   string
	port int
}

type NetConf struct {
	TestConfData struct {
		Ip   string `json:"ip"`
		Port int    `json:"port"`
	} `json:"testConfData"`
	TestName string `json:"TestName"`
}

var instance *ConfMgr
var netConf *NetConf

func GetInstance() *ConfMgr {
	if instance == nil {
		instance = &ConfMgr{}
	}
	return instance
}

func Init(filePath string) {
	confFile, err := os.Open(filePath)
	defer confFile.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	json.NewDecoder(confFile).Decode(netConf)
}
