package vm

import "log"

type Operation interface {
	Execute(ef *ExecutionFrame)
}

type PopToRegisterOperation struct {
	registerType RegisterType
	index        int
}

func (o *PopToRegisterOperation) Execute(ef *ExecutionFrame) {
	log.Println("Popping to register", o.index)
	switch o.registerType {
	case StringRegisterType:
		val := ef.valueStack.pop()
		if s, ok := val.(*StringValue); ok {
			ef.stringRegisters[o.index] = *s
			log.Println("Popped to register", o.index, s.Value)
		} else {
			log.Panicln("Value is not a string")
		}
	}
}

type PushContextOperation struct {
	contextNameIndex int
}

func (o *PushContextOperation) Execute(ef *ExecutionFrame) {
	log.Println("Pushing context", ef.GetFromStringPool(o.contextNameIndex))
	contextName := ef.GetFromStringPool(o.contextNameIndex)
	context := ef.GetObjectFromContext(contextName)
	ef.valueStack.push(context)
	log.Println("Pushed context", contextName)
}

type MethodCallOperation struct {
	argumentCount int
}

func (o *MethodCallOperation) Execute(ef *ExecutionFrame) {
	log.Println("Calling")
	var object = ef.valueStack.pop()

	objectValue, ok := object.(ObjectValue)
	if !ok {
		log.Panicln("Value is not an object")
	}

	var method = ef.valueStack.pop()

	if _, ok := method.(*StringValue); !ok {
		log.Panicln("Value is not a method")
	}

	ef.call(objectValue, method)
	log.Println("Called", object, method)
}

type AddOperation struct{}

func (o *AddOperation) Execute(ef *ExecutionFrame) {
	log.Println("Adding")
	var a = ef.valueStack.pop()
	var b = ef.valueStack.pop()
	c := b.AddValue(a)

	ef.valueStack.push(c)
	log.Println("Added", a, b)
	log.Println("Result", c)
}

type PushStringOperation struct {
	index int
}

func (o *PushStringOperation) Execute(ef *ExecutionFrame) {
	log.Println("Pushing string", ef.GetFromStringPool(o.index))
	ef.valueStack.push(NewStringValue(ef.GetFromStringPool(o.index)))
	log.Println("Pushed string", ef.GetFromStringPool(o.index))
}
