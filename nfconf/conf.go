package nfconf

import (
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
			e = err
			return
		}

		Conf, err = simplejson.NewJson(rf)
		if err != nil {
			e = err
			return
		}
		e = err
	}
	once.Do(initFunc)
	return e
}
