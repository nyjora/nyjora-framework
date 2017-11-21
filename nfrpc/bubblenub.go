
package nfrpc

import (
	"fmt"

	"nyjora-framework/nflog"
)

type BubbleNub struct {
	local NubInfo
	dispatcher Dispatcher
	emitter Emitter
	forwarder NubForwarder
}

func (n *BubbleNub) HandleBubble(remote NubInfo, methodid MethodID,
	data []byte) error {
	if methodid > 0 { // request or message, to dispatcher
		if n.dispatcher == nil {
			nflog.Err("Nub.HandleBubble with dispatcher == nil")
			return fmt.Errorf("HandleBubble with dispatcher == nil")
		}
		n.dispatcher.HandleBubble(remote, methodid, data)
	} else {
		if n.emitter == nil {
			nflog.Err("Nub.HandleBubble with emitter == nil")
			return fmt.Errorf("HandleBubble with emitter == nil")
		}
	}
	return nil
}

func (n *BubbleNub) ForwardBubble(remote NubInfo,
	methodid MethodID, data []byte) error {
	if n.forwarder == nil {
		nflog.Err("Nub.ForwardBubble with forwarder == nil")
		return fmt.Errorf("ForwardBubble with forwarder == nil")
	}
	return n.forwarder.ForwardBubble(remote, n.local, methodid, data)
}