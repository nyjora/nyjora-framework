package nfconf

type ConfMgr struct {
}

var instance *ConfMgr

func GetInstance() *ConfMgr {
	if instance == nil {
		instance = &ConfMgr{}
	}
	return instance
}
