package nfrpc

type MethodID = int32
type ReplyID = uint64

type Emitter interface {
	HandleBubble(remote NubInfo, methodid MethodID, data []byte) error
}

type Dispatcher interface {
	HandleBubble(remote NubInfo, methodid MethodID, data []byte) error
}

type NubForwarder interface {
	ForwardBubble(remote NubInfo, local NubInfo,
		methodid MethodID, data []byte) error
}