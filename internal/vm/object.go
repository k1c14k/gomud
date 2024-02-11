package vm

type Object interface {
	GetClass() Class
}

type vmObject struct {
	class Class
}

func (o *vmObject) GetClass() Class {
	return o.class
}

func NewObject(name string) Object {
	class := instance.getClass(name)
	return &vmObject{class}
}
