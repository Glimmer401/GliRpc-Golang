package service

type Procedure struct {
	Name string			// Remote Procedure name
	Args interface{}	// Remote Procedure argument
	Rets interface{}	// Remote Procedure return values
	Index uint32		// an index to manage
}

type Service struct {
	Name2Index map[string]uint32
	Procedures Procedure[]
}

func Register(handler func) {
	 
}