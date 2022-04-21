package service

import (
	"reflect"
	"sync/atomic"
)

type MethodType struct {
	method		reflect.Method
	ArgsType 	reflect.Type
	ReplyType 	reflect.Type
	numCalls	uint64			// counter for a certain method
}

func (m *MethodType) NumCalls() uint64 {
	return atomic.LoadUint64(&m.numCalls)
}

func (m *MethodType) NewArgv() reflect.Value {
	var argv reflect.Value
	if m.ArgsType.Kind() == reflect.Ptr {
		argv = reflect.New(m.ArgsType.Elem())
	} else {
		argv = reflect.New(m.ArgsType).Elem()
	}
	return argv
}

func (m *MethodType) NewReplyv() reflect.Value {
	var replyv reflect.Value
	replyv = reflect.New(m.ReplyType.Elem())
	switch m.ReplyType.Kind() {
	case reflect.Map:
		replyv.Elem().Set(reflect.MakeMap(m.ReplyType.Elem()))
	case reflect.Slice:
		replyv.Elem().Set(reflect.MakeSlice(m.ReplyType.Elem(), 0, 0))
	}
	return replyv
}
