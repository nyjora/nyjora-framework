package nfconf

import (
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/bitly/go-simplejson"
)

var Conf *simplejson.Json
var once sync.Once

func Init(filePath string) (e error) {
	var initFunc = func() {
		rf, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println(err.Error())
			e = err
			return
		}

		Conf, err = simplejson.NewJson(rf)
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
