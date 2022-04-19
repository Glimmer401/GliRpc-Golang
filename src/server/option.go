package server

import (
	"erpc/codec"
)

const Magic = 0x0817ff
type Option struct {
	Magic		int        // Magic marks this's a erpc request
	CodecType   codec.Type // different Codec to encode body
}

var DefaultOption = &Option{
	Magic: 		Magic,
	CodecType:  codec.GobType,
}