package main

import (
	"erpc/service"
	"erpc/example"
	"erpc/util"

	_ "reflect"
	_ "fmt"
)

func main() {
	// start a server
	service := service.Service{}
	service.Init("testing", 10)
	service.Register(example.Handler, example.ExamArgs{1,1}, example.ExamRets{0})

	// transfer a raw request to request object
	rawRequest := "{\"name\":\"main.handler\",\"args\":{\"X\":1,\"Y\":1},\"Rets\":{\"Z\":0}}"	
	request := util.Request{}
	request.Decode(rawRequest)
	

	// apply the procedure
	// fmt.Println(service.Apply(request.Name, &args, &rets))
	service.Apply(request)
}
