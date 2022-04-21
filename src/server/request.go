package server

import (
	"erpc/codec"
	"erpc/service"
	"reflect"
)

// request stores all information of a call
type request struct {
	h            *codec.Header // header of request
	argv, replyv reflect.Value // argv and replyv of request
	svc			 *service.Service	
	mtype		 *service.MethodType
}

var invalidRequest = struct{}{}