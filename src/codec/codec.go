package codec

import "io"

type Header struct {
	MethodName string		// format "Service.Method"
	Error 	   string
	Seq		   uint64
}

type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

type NewCodecFunc func(io.ReadWriteCloser) Codec
type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json" 
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType]  = NewGobCodec
	NewCodecFuncMap[JsonType] = NewJsonCodec
}