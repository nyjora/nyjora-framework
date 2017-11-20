
package nfrpc

import (
	"nyjora-framework/nflog"
)

type BubbleNub struct {
	dispatcher Dispatcher
	emitter Emitter
}

func (n *BubbleNub) HandleBubble(remote RemoteNubInfo, methodid MethodID,
	data *bytes.Buffer) error {
	if methodid > 0 { // request or message, to dispatcher
		if n.dispatcher == nil {
			
		}
	}
}
