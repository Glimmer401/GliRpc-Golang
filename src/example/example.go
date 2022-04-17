package example

// example RPC args and Reture value type
type ExamArgs struct {
	x int
	y int
}
type ExamReture struct  {
	z int
}

// example target remote procedure call
func handler(int x, int y) int {  return x+y  }