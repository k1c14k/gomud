package vm

import "bytes"

type Object struct {
	class Class
}

func (o *Object) GetClass() Class {
	return o.class
}

func (o *Object) String() string {
	buff := bytes.Buffer{}
	buff.WriteString("Object[")
	buff.WriteString(o.class.String())
	buff.WriteString("]")
	return buff.String()
}

func NewObject(name string) *Object {
	class := instance.getClass(name)
	return &Object{class}
}

func NewObjectFromClass(class Class) *Object {
	return &Object{class}
}
