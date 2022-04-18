package example

// example RPC args ,Reture value type, and a simple handler
type ExamArgs struct {
	x int
	y int
}

type ExamRets struct  {
	z int
}

func handler(args *ExamArgs, rets *ExamRets) {
	rets.z = args.x + args.y
}