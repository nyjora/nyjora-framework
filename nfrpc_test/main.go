
package nfrpc_test

import (
	"nfrpc"
)
//go:generate protoc --go_out=. nfrpc_test.proto
//go:generate ../nfrpc/nfrpc-gen.py --go_out=. nfrpc_test.xml
