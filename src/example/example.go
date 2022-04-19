package example

// example RPC args ,Reture value type, and a simple handler
type ExamArgs struct {
	X int	`json:"X`
	Y int	`json:"Y`
}

type ExamRets struct  {
	Z int	`json:"Z`
}

func Handler(args *ExamArgs, rets *ExamRets) {
	rets.Z = args.X + args.Y
}