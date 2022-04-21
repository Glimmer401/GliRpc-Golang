package service

import (
	"log"
	"reflect"
	"go/ast"
	"sync/atomic"
)

type Service struct {
	Name		string
	receiver	reflect.Value
	typ			reflect.Type
	Methods		map[string]*MethodType
}


// create a new service from a struct
func NewService (rcvr interface{}) *Service {
	s := new(Service)
	s.receiver = reflect.ValueOf(rcvr)
	s.typ      = reflect.TypeOf(rcvr)
	s.Name     = reflect.Indirect(s.receiver).Type().Name()

	if !ast.IsExported(s.Name) {
		log.Fatalf("erpc server: %s is not a exported service name", s.Name)
	}

	s.registerMethods()
	return s
}


// after init a service, register its methods
func (s *Service) registerMethods() {
	s.Methods = make(map[string]*MethodType)
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
		if !isBuiltinOrIsExported(argType) || !isBuiltinOrIsExported(replyType) {
			continue
		}
		s.Methods[method.Name] = &MethodType{
			method:    method,
			ArgsType:   argType,
			ReplyType: replyType,
		}
		log.Printf("erpc server: register %s.%s\n", s.Name, method.Name)
	}
}

// check a type if it is a builtin or exported type
func isBuiltinOrIsExported(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}


// call a method
func (s *Service) Call(m *MethodType, argv, replyv reflect.Value) error {
	atomic.AddUint64(&m.numCalls, 1)
	f := m.method.Func
	returnValues := f.Call([]reflect.Value{s.receiver, argv, replyv})
	if errInter := returnValues[0].Interface(); errInter != nil {
		return errInter.(error)
	}
	return nil
}


