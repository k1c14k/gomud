package vm

import "log"

type RegisterType int

const (
	StringRegisterType RegisterType = iota
)

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

type ExecutionContext struct {
	stringPool    []string
	objectContext map[string]ObjectValue
}

func NewExecutionContext(stringPool []string, objectContext map[string]ObjectValue) ExecutionContext {
	return ExecutionContext{
		stringPool:    stringPool,
		objectContext: objectContext,
	}
}

type ExecutionFrame struct {
	stringRegisters  []StringValue
	valueStack       ValueStack
	executionContext ExecutionContext
	nextFrame        *ExecutionFrame
	programCounter   int
	program          []Operation
}

func NewExecutionFrame(ctx ExecutionContext) *ExecutionFrame {
	return &ExecutionFrame{
		stringRegisters:  make([]StringValue, 20),
		valueStack:       *NewValueStack(),
		executionContext: ctx,
		programCounter:   0,
	}
}

func (ef *ExecutionFrame) GetFromStringPool(index int) string {
	if index >= len(ef.executionContext.stringPool) {
		log.Panicln("String pool index out of range")
	}
	return ef.executionContext.stringPool[index]
}

func (ef *ExecutionFrame) GetObjectFromContext(name string) ObjectValue {
	if obj, ok := ef.executionContext.objectContext[name]; ok {
		return obj
	}
	log.Panicln("Object not found in context")
	return ObjectValue{}
}

func (ef *ExecutionFrame) call(object ObjectValue, method Value, arguments []Value) {
	objectValue := *object.value
	ef.nextFrame = NewExecutionFrame(ef.executionContext)
	cls := objectValue.GetClass()
	m := cls.GetMethod(method.(*StringValue).Value)
	for _, a := range arguments {
		ef.nextFrame.valueStack.push(a)
	}
	switch objectValue.(type) {
	case *vmObject:
		ef.nextFrame.program = m.(*vmMethod).operations
		ef.nextFrame.run()
	default:
		m.Execute(ef.nextFrame)
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
