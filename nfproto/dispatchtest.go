package nfproto

import (
	"fmt"
	"nyjora-framework/nfcommon"
	"nyjora-framework/nfnet"
	"nyjora-framework/nftest/nflogic_test"
)

func DispatchTest(proto *nfcommon.Protocol, session *nfnet.NetSession) {
	if proto.Id == 1 {
		// c2s
		nflogic_test.OnDispatch(session.RemoteAddr(), proto)
	} else if proto.Id == 2 {
		// s2c
		fmt.Printf("says : %s\n", proto.Data)
	}
}

func Init() {
	nfnet.TestRpc = DispatchTest
}
