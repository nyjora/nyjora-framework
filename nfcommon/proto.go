package nfcommon

type Protocol struct {
	Id       uint32
	FromType uint32
	FromId   uint32
	ToType   uint32
	ToId     uint32
	Data     []byte
}

func NewProto() *Protocol {
	return &Protocol{}
}
