package vm

import "log"

type ValueStack struct {
	values  []Value
	pos     int
	maxSize int
}

func NewValueStack() *ValueStack {
	maxSize := 20
	return &ValueStack{
		values:  make([]Value, maxSize),
		maxSize: maxSize,
	}
}

func (vs *ValueStack) pop() Value {
	if vs.pos == 0 {
		log.Panicln("ValueStack is empty")
	}
	vs.pos--
	return vs.values[vs.pos]
}

func (vs *ValueStack) push(v Value) {
	if vs.pos == vs.maxSize {
		log.Panicln("ValueStack is full")
	}
	vs.values[vs.pos] = v
	vs.pos++
}
