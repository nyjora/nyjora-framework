
package nfrpc_test

import (
	_ "nyjora-framework/nfrpc"
)
//go:generate protoc --go_out=. nfrpc_test.proto
//go:generate ../nfrpc/nfrpc-gen.py --go_out=. nfrpc_test.xml
