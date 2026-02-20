package vm

import (
	"strconv"

	"github.com/sirupsen/logrus"
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
	logrus.Debug("Popping to register", o.index)
	ef.registers[o.index] = ef.valueStack.pop()
}

func (o *PopToRegisterOperation) String() string {
	return "RPOP " + strconv.Itoa(int(o.registerType)) + " " + strconv.Itoa(o.index)
}

type PushContextOperation struct {
	contextNameIndex int
}

func (o *PushContextOperation) Execute(ef *ExecutionFrame) {
	logrus.Debug("Pushing context", ef.GetFromStringPool(o.contextNameIndex))
	contextName := ef.GetFromStringPool(o.contextNameIndex)
	context := ef.GetObjectFromContext(contextName)
	ef.valueStack.push(context)
	logrus.Debug("Pushed context", contextName)
}

func (o *PushContextOperation) String() string {
	return "CPUSH " + strconv.Itoa(o.contextNameIndex)
}

type MethodCallOperation struct {
}

func (o *MethodCallOperation) Execute(ef *ExecutionFrame) {
	logrus.Debug("Calling")
	var object = ef.valueStack.pop()

	objectValue, ok := object.(ObjectValue)
	if !ok {
		logrus.Panic("Value is not an object")
	}

	var method = ef.valueStack.pop()

	if _, ok := method.(*StringValue); !ok {
		logrus.Panic("Value is not a method")
	}

	ef.call(objectValue, method)
	logrus.Debug("Called", object, method)
}

func (o *MethodCallOperation) String() string {
	return "CALL"
}

type AddOperation struct{}

func (o *AddOperation) Execute(ef *ExecutionFrame) {
	logrus.Debug("Adding")
	var a = ef.valueStack.pop()
	var b = ef.valueStack.pop()
	c := b.Add(a)

	ef.valueStack.push(c)
	logrus.Debug("Added", a, b)
	logrus.Debug("Result", c)
}

func (o *AddOperation) String() string {
	return "ADD"
}

type SubOperation struct{}

func (o *SubOperation) String() string {
	return "SUB"
}

func (o *SubOperation) Execute(ef *ExecutionFrame) {
	logrus.Debug("Subtracting")
	var a = ef.valueStack.pop()
	var b = ef.valueStack.pop()
	c := b.Subtract(a)

	ef.valueStack.push(c)
	logrus.Debug("Subtracted", a, b)
	logrus.Debug("Result", c)

}

type MulOperation struct{}

func (o *MulOperation) String() string {
	return "MUL"
}

func (o *MulOperation) Execute(ef *ExecutionFrame) {
	logrus.Debug("Multiplying")
	var a = ef.valueStack.pop()
	var b = ef.valueStack.pop()
	c := b.Multiply(a)

	ef.valueStack.push(c)
	logrus.Debug("Multiplied", a, b)
	logrus.Debug("Result", c)
}

type DivOperation struct{}

func (o *DivOperation) String() string {
	return "DIV"
}

func (o *DivOperation) Execute(ef *ExecutionFrame) {
	logrus.Debug("Dividing")
	var a = ef.valueStack.pop()
	var b = ef.valueStack.pop()
	c := b.Divide(a)

	ef.valueStack.push(c)
	logrus.Debug("Divided", a, b)
	logrus.Debug("Result", c)
}

type ModOperation struct{}

func (m ModOperation) Execute(ef *ExecutionFrame) {
	logrus.Debug("Modding")
	var a = ef.valueStack.pop()
	var b = ef.valueStack.pop()
	c := b.Modulo(a)

	ef.valueStack.push(c)
	logrus.Debug("Modded", a, b)
	logrus.Debug("Result", c)
}

func (m ModOperation) String() string {
	return "MOD"
}

type PushStringOperation struct {
	index int
}

func (o *PushStringOperation) Execute(ef *ExecutionFrame) {
	logrus.Debug("Pushing string", ef.GetFromStringPool(o.index))
	ef.valueStack.push(NewStringValue(ef.GetFromStringPool(o.index)))
	logrus.Debug("Pushed string", ef.GetFromStringPool(o.index))
}

func (o *PushStringOperation) String() string {
	return "SPUSH " + strconv.Itoa(o.index)
}

type JumpIfFalseOperation struct {
	target int
}

func (o *JumpIfFalseOperation) Execute(ef *ExecutionFrame) {
	logrus.Debug("Jumping if false")
	var a = ef.valueStack.pop()
	if !a.isTruthy() {
		ef.programCounter = o.target - 1
	}
	logrus.Debug("Jumped if false", a)
}

func (o *JumpIfFalseOperation) String() string {
	return "JMPF " + strconv.Itoa(o.target)
}

type JumpOperation struct {
	target int
}

func (o *JumpOperation) Execute(ef *ExecutionFrame) {
	logrus.Debug("Jumping")
	ef.programCounter = o.target - 1
	logrus.Debug("Jumped")
}

func (o *JumpOperation) String() string {
	return "JMP " + strconv.Itoa(o.target)
}

type EqualOperation struct{}

func (o *EqualOperation) Execute(ef *ExecutionFrame) {
	logrus.Debug("Comparing")
	var a = ef.valueStack.pop()
	var b = ef.valueStack.pop()
	c := a.equalValue(b)
	ef.valueStack.push(c)
	logrus.Debug("Compared", a, b)
	logrus.Debug("Result", c)
}

func (o *EqualOperation) String() string {
	return "EQ"
}

type PushFromRegisterOperation struct {
	registerType RegisterType
	index        int
}

func (o *PushFromRegisterOperation) Execute(ef *ExecutionFrame) {
	logrus.Debug("Pushing from register", o.index)
	ef.valueStack.push(ef.registers[o.index])
	logrus.Debug("Pushed from register", o.index)
}

func (o *PushFromRegisterOperation) String() string {
	return "RPUSH " + strconv.Itoa(int(o.registerType)) + strconv.Itoa(o.index)
}

type ReturnOperation struct{}

func (o *ReturnOperation) Execute(ef *ExecutionFrame) {
	logrus.Debug("Returning")
	ef.programCounter = len(ef.program)
	logrus.Debug("Returned")
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
	logrus.Debug("Pushing number", o.value)
	ef.valueStack.push(NewNumberValue(o.value))
	logrus.Debug("Pushed number", o.value)
}
