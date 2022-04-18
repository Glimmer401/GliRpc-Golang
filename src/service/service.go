package service

import (
	"fmt"
	"reflect"
	"runtime"
)

// an request for procedure
type Procedure struct {
	Name string			// Remote Procedure Func name
	Handler interface{} // Remote Procedure Handlers
	Args interface{}	// Remote Procedure Args value
	Rets interface{} 	// Remote Procedure Rets value
}

// registered procedures
type Service struct {
	Name 	   string			// service name
	Handlers   []Procedure	    // procedures lists
	Name2Index map[string]int	// name to lists index
}

// init Service with a capacity
func (s *Service) Init(name string, capacity int) {
	s.Name = name
	s.Handler = make([]Procedure, 0, capacity)
	s.Name2Index = make(map[string]int)
}

// register a procedure into services 
func (s *Service) Register(function interface{}) {
	if len(s.Func) ==  cap(s.Func) {
		fmt.Printf("%s service is full loaded and can not register any more procedure", s.Name)
		return
	}
	fn := runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()
	s.Name2Index[fn] = len(s.Func)
	s.Func = append(s.Func, (function))
}

// apply a registered procedure
func (s *Service) Apply(handler string, args interface{}, rets interface{}) ([]reflect.Value) {
	procedure := reflect.ValueOf(s.Func[s.Name2Index[handler]])
	input := make([]reflect.Value, 2)
	input[0] = reflect.ValueOf(args)
	input[1] = reflect.ValueOf(rets)
	resp := procedure.Call(input)
	return resp
}