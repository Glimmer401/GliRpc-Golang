package server

import (
	"erpc/codec"
	"reflect"
)

// request stores all information of a call
type request struct {
	h            *codec.Header // header of request
	argv, replyv reflect.Value // argv and replyv of request
}

var invalidRequest = struct{}{}