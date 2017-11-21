
package nfrpc

type NubInfo interface {
	remote_type() int32
	remote_id() int64
}