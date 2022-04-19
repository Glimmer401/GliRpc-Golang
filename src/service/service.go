package service

import (
	"fmt"
	"reflect"
	"runtime"
	"erpc/util"
	"encoding/json"
)

// an request for procedure
type Procedure struct {
	Name string			// Remote Procedure Func name
	Handler interface{} // Remote Procedure Handlers
	Args    reflect.Type// Remote Procedure Args type
	Rets 	reflect.Type// Remote Procedure Rets type
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
	s.Handlers = make([]Procedure, 0, capacity)
	s.Name2Index = make(map[string]int)
}

// register a procedure into services 
func (s *Service) Register(function interface{}, args interface{}, rets interface{}) {
	if len(s.Handlers) ==  cap(s.Handlers) {
		fmt.Printf("%s service is full loaded and can not register any more procedure", s.Name)
		return
	}
	fn := runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()
	s.Name2Index[fn] = len(s.Handlers)
	s.Handlers = append(s.Handlers, 
						(Procedure{fn, function, reflect.TypeOf(args), reflect.TypeOf(rets)}))
}

// apply a registered procedure
// input should be a certain object
func (s *Service) Apply(request util.Request) ([]reflect.Value) {
	procedure := s.Handlers[s.Name2Index[request.Name]]
	handler := reflect.ValueOf(procedure.Handler)

	request.Args = request.Args.(map[string]interface{})
	request.Rets = request.Rets.(map[string]interface{})
	argsJson, _ := json.Marshal(request.Args)
	retsJson, _ := json.Marshal(request.Rets)
	fmt.Println(string(argsJson))
	args := reflect.New(procedure.Args)
	rets := reflect.New(procedure.Rets)
	json.Unmarshal(argsJson, args)
	json.Unmarshal(retsJson, rets)

	fmt.Println(args)
	fmt.Println(rets)

	json.Unmarshal(argsJson, &args)
	json.Unmarshal(retsJson, &rets)
	
	input := make([]reflect.Value, 2)
	input[0] = args
	input[1] = rets

	
	handler.Call(input)
	
	return nil
}