package example

// example RPC args ,Reture value type, and a simple handler
type ExamArgs struct {
	X int	`json:"X`
	Y int	`json:"Y`
}

type ExamRets struct  {
	Z int	`json:"Z`
}

type Calc struct {

}

func (c *Calc) Add(args ExamArgs, rets *ExamRets) error{
	rets.Z = args.X + args.Y
	return nil
}