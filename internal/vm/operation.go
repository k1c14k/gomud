package vm

import (
	"log"
	"strconv"
)

type Operation interface {
	Execute(ef *ExecutionFrame)
	String() string
}

type PopToRegisterOperation struct {
	registerType RegisterType
	index        int
}

func (o *PopToRegisterOperation) Execute(ef *ExecutionFrame) {
	log.Println("Popping to register", o.index)
	ef.registers[o.index] = ef.valueStack.pop()
}

func (o *PopToRegisterOperation) String() string {
	return "RPOP " + strconv.Itoa(int(o.registerType)) + " " + strconv.Itoa(o.index)
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

func (o *PushContextOperation) String() string {
	return "CPUSH " + strconv.Itoa(o.contextNameIndex)
}

type MethodCallOperation struct {
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

func (o *MethodCallOperation) String() string {
	return "CALL"
}

type AddOperation struct{}

func (o *AddOperation) Execute(ef *ExecutionFrame) {
	log.Println("Adding")
	var a = ef.valueStack.pop()
	var b = ef.valueStack.pop()
	c := b.Add(a)

	ef.valueStack.push(c)
	log.Println("Added", a, b)
	log.Println("Result", c)
}

func (o *AddOperation) String() string {
	return "ADD"
}

type SubOperation struct{}

func (o *SubOperation) String() string {
	return "SUB"
}

func (o *SubOperation) Execute(ef *ExecutionFrame) {
	log.Println("Subtracting")
	var a = ef.valueStack.pop()
	var b = ef.valueStack.pop()
	c := b.Subtract(a)

	ef.valueStack.push(c)
	log.Println("Subtracted", a, b)
	log.Println("Result", c)

}

type MulOperation struct{}

func (o *MulOperation) String() string {
	return "MUL"
}

func (o *MulOperation) Execute(ef *ExecutionFrame) {
	log.Println("Multiplying")
	var a = ef.valueStack.pop()
	var b = ef.valueStack.pop()
	c := b.Multiply(a)

	ef.valueStack.push(c)
	log.Println("Multiplied", a, b)
	log.Println("Result", c)
}

type DivOperation struct{}

func (o *DivOperation) String() string {
	return "DIV"
}

func (o *DivOperation) Execute(ef *ExecutionFrame) {
	log.Println("Dividing")
	var a = ef.valueStack.pop()
	var b = ef.valueStack.pop()
	c := b.Divide(a)

	ef.valueStack.push(c)
	log.Println("Divided", a, b)
	log.Println("Result", c)
}

type ModOperation struct{}

func (m ModOperation) Execute(ef *ExecutionFrame) {
	log.Println("Modding")
	var a = ef.valueStack.pop()
	var b = ef.valueStack.pop()
	c := b.Modulo(a)

	ef.valueStack.push(c)
	log.Println("Modded", a, b)
	log.Println("Result", c)
}

func (m ModOperation) String() string {
	return "MOD"
}

type PushStringOperation struct {
	index int
}

func (o *PushStringOperation) Execute(ef *ExecutionFrame) {
	log.Println("Pushing string", ef.GetFromStringPool(o.index))
	ef.valueStack.push(NewStringValue(ef.GetFromStringPool(o.index)))
	log.Println("Pushed string", ef.GetFromStringPool(o.index))
}

func (o *PushStringOperation) String() string {
	return "SPUSH " + strconv.Itoa(o.index)
}

type JumpIfFalseOperation struct {
	target int
}

func (o *JumpIfFalseOperation) Execute(ef *ExecutionFrame) {
	log.Println("Jumping if false")
	var a = ef.valueStack.pop()
	if !a.isTruthy() {
		ef.programCounter = o.target - 1
	}
	log.Println("Jumped if false", a)
}

func (o *JumpIfFalseOperation) String() string {
	return "JMPF " + strconv.Itoa(o.target)
}

type JumpOperation struct {
	target int
}

func (o *JumpOperation) Execute(ef *ExecutionFrame) {
	log.Println("Jumping")
	ef.programCounter = o.target - 1
	log.Println("Jumped")
}

func (o *JumpOperation) String() string {
	return "JMP " + strconv.Itoa(o.target)
}

type EqualOperation struct{}

func (o *EqualOperation) Execute(ef *ExecutionFrame) {
	log.Println("Comparing")
	var a = ef.valueStack.pop()
	var b = ef.valueStack.pop()
	c := a.equalValue(b)
	ef.valueStack.push(c)
	log.Println("Compared", a, b)
	log.Println("Result", c)
}

func (o *EqualOperation) String() string {
	return "EQ"
}

type PushFromRegisterOperation struct {
	registerType RegisterType
	index        int
}

func (o *PushFromRegisterOperation) Execute(ef *ExecutionFrame) {
	log.Println("Pushing from register", o.index)
	ef.valueStack.push(ef.registers[o.index])
	log.Println("Pushed from register", o.index)
}

func (o *PushFromRegisterOperation) String() string {
	return "RPUSH " + strconv.Itoa(int(o.registerType)) + strconv.Itoa(o.index)
}

type ReturnOperation struct{}

func (o *ReturnOperation) Execute(ef *ExecutionFrame) {
	log.Println("Returning")
	ef.programCounter = len(ef.program)
	log.Println("Returned")
}

func (o *ReturnOperation) String() string {
	return "RET"
}

type PushNumberOperation struct {
	value int
}

func (o *PushNumberOperation) String() string {
	return "PUSN " + strconv.Itoa(o.value)
}

func (o *PushNumberOperation) Execute(ef *ExecutionFrame) {
	log.Println("Pushing number", o.value)
	ef.valueStack.push(NewNumberValue(o.value))
	log.Println("Pushed number", o.value)
}
