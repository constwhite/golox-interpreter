package interpreter

import "time"

type Clock struct {
}

func (c Clock) arity() int {
	return 0
}

func (c Clock) call(interpreter *Interpreter, arguements []interface{}) interface{} {
	return float64(time.Now().UnixMilli() / 1000)
}
