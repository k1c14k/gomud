package vm

import "log"

type RegisterType int

const (
	StringRegisterType RegisterType = iota
)

type ContextProvider interface {
	GetObjectValueFromContext(name string) *ObjectValue
}

type ExecutionFrame struct {
	registers       []Value
	valueStack      ValueStack
	nextFrame       *ExecutionFrame
	programCounter  int
	program         []Operation
	stringPool      []string
	contextProvider ContextProvider
}

func NewExecutionFrame(contextProvider ContextProvider) *ExecutionFrame {
	return &ExecutionFrame{
		registers:       make([]Value, 20),
		valueStack:      *NewValueStack(),
		programCounter:  0,
		contextProvider: contextProvider,
	}
}

func (ef *ExecutionFrame) GetFromStringPool(index int) string {
	if index >= len(ef.stringPool) {
		log.Panicln("String pool index out of range")
	}
	return ef.stringPool[index]
}

func (ef *ExecutionFrame) GetObjectFromContext(name string) ObjectValue {
	obj := ef.contextProvider.GetObjectValueFromContext(name)
	if obj == nil {
		log.Panicln("Object not found in context")
	}
	return *obj
}

func (ef *ExecutionFrame) call(object ObjectValue, method Value) {
	objectValue := *object.value
	cls := objectValue.GetClass()
	m := cls.GetMethod(method.(*StringValue).Value)
	switch m.(type) {
	case *vmMethod:
		ef.nextFrame = NewExecutionFrame(ef.contextProvider)
		for i := m.GetArgumentCount() - 1; i >= 0; i-- {
			ef.nextFrame.valueStack.push(ef.valueStack.pop())
		}
		ef.nextFrame.program = m.(*vmMethod).operations
		ef.nextFrame.stringPool = m.(*vmMethod).GetStrings()
		ef.nextFrame.run()
		for i := 0; i < m.GetReturnValueCount(); i++ {
			ef.valueStack.push(ef.nextFrame.valueStack.pop())
		}
		ef.nextFrame = nil
	case *internalMethod:
		arguments := make([]Value, m.GetArgumentCount())
		for i := m.GetArgumentCount() - 1; i >= 0; i-- {
			arguments[i] = ef.valueStack.pop()
		}
		result := m.(*internalMethod).handle(arguments)
		for _, r := range result {
			ef.valueStack.push(r)
		}
	}
}

func (ef *ExecutionFrame) PopValue() Value {
	return ef.valueStack.pop()
}

func (ef *ExecutionFrame) run() {
	for ef.programCounter < len(ef.program) {
		ef.program[ef.programCounter].Execute(ef)
		ef.programCounter++
	}
}
