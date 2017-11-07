package nflogic_test

import (
	"fmt"
	"net"
	"nyjora-framework/nfcommon"
	"nyjora-framework/nftest/nfservice_test"
)

var Nftestserver *nfservice_test.TcpServerDelegate

func OnDispatch(addr net.Addr, proto *nfcommon.Protocol) {
	fmt.Printf("%s:%d,%d,%d,%d,%d,%s\n", addr, proto.Id, proto.FromType, proto.FromId, proto.ToType, proto.ToId, proto.Data)
	Nftestserver.BroadCastChatTest(2, 3, 4, 5, 6, proto.Data)
}
