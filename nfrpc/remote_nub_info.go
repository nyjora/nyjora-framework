
package nfrpc

type RemoteNubInfo interface {
	remote_type() int32
	remote_id() int64
}