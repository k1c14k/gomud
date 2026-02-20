package vm

import "github.com/sirupsen/logrus"

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
		logrus.Panic("ValueStack is empty")
	}
	vs.pos--
	return vs.values[vs.pos]
}

func (vs *ValueStack) push(v Value) {
	if vs.pos == vs.maxSize {
		logrus.Panic("ValueStack is full")
	}
	vs.values[vs.pos] = v
	vs.pos++
}
