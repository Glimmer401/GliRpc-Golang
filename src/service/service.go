package service

import (
	"log"
	"reflect"
	"go/ast"
	"sync/atomic"
)

type service struct {
	name		string
	receiver	reflect.Value
	typ			reflect.Type
	methods		map[string]*methodType
}


// create a new service from a struct
func newService (rcvr interface{}) *service {
	s := new(service)
	s.receiver = reflect.ValueOf(rcvr)
	s.typ      = reflect.TypeOf(rcvr)
	s.name     = reflect.Indirect(s.receiver).Type().Name()

	if !ast.IsExported(s.name) {
		log.Fatalf("erpc server: %s is not a exported service name", s.name)
	}

	s.registerMethods()
	return s
}


// after init a service, register its methods
func (s *service) registerMethods() {
	s.method = make(map[string]*methodType)
	for i := 0; i < s.typ.NumMethod(); i++ {
		method := s.typ.Method(i)
		mType := method.Type
		if mType.NumIn() != 3 || mType.NumOut() != 1 {
			continue
		}
		if mType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}
		argType, replyType := mType.In(1), mType.In(2)
		if !isExportedOrBuiltinType(argType) || !isExportedOrBuiltinType(replyType) {
			continue
		}
		s.method[method.Name] = &methodType{
			method:    method,
			ArgType:   argType,
			ReplyType: replyType,
		}
		log.Printf("erpc server: register %s.%s\n", s.name, method.Name)
	}
}

// check a type if it is a builtin or exported type
func isBuiltinOrIsExported(t reflect.Value) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}


// call a method
func (s *service) call(m *methodType, argv, replyv reflect.Value) error {
	atomic.AddUint64(&m.numCalls, 1)
	f := m.method.Func
	returnValues := f.Call([]reflect.Value{s.rcvr, argv, replyv})
	if errInter := returnValues[0].Interface(); errInter != nil {
		return errInter.(error)
	}
	return nil
}


