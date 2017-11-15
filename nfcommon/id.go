package nfcommon

type SessionId int64

var curSessionId SessionId = 0

func NextSessionId() SessionId {
	curSessionId++
	return curSessionId
}

func (id SessionId) IsNil() bool {
	return id <= 0
}
