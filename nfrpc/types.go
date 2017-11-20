package nfrpc

import(
	"bytes"
)

type MethodID = uint32
type ReplyID = uint64

type Emitter interface {
	HandleBubble(remote RemoteNubInfo, methodid MethodID, data *bytes.Buffer)
}

type Dispatcher interface {
	HandleBubble(remote RemoteNubInfo, methodid MethodID, data *bytes.Buffer)
}

