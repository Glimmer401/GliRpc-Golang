package client

// describe a remote call
type Call struct{
	// index
	Seq				uint64
	// info for server
	MethodName		string
	Args			interface{}
	Reply			interface{}
	// status
	Error			error
	Done			chan *Call
}

func (call *Call) done() {
	call.Done <- call
}