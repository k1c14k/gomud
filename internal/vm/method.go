package vm

import (
	"goMud/internal/gmsl/compiler"
)

type Method interface {
	GetArgumentCount() int
	GetReturnValueCount() int
}

type vmMethod struct {
	argumentCount    int
	returnValueCount int
	operations       []Operation
	strings          []string
}

type MethodHandler func([]Value) []Value

type internalMethod struct {
	argumentCount    int
	returnValueCount int
	handle           MethodHandler
}

func (m *internalMethod) GetArgumentCount() int {
	return m.argumentCount
}

func (m *internalMethod) GetReturnValueCount() int {
	return m.returnValueCount
}

func (m *vmMethod) GetArgumentCount() int {
	return m.argumentCount
}

func (m *vmMethod) GetReturnValueCount() int {
	return m.returnValueCount
}

func NewMethodsFromAssembly(aOut *compiler.Assembly) map[string]Method {
	result := make(map[string]Method)

	for _, f := range aOut.GetFunctions() {
		name, method := f.GetName(), NewMethodFromAssembly(f)
		result[name] = method

	}

	return result
}

func NewMethodFromAssembly(f compiler.FunctionInfo) Method {
	result := &vmMethod{argumentCount: f.GetArgumentCount(), returnValueCount: f.GetReturnValueCount(), strings: f.GetStrings(), operations: make([]Operation, 0)}
	labelPos := make(map[string]int)
	posRequestingLabel := make(map[int]string)

	for _, e := range f.GetEntries() {
		switch e.(type) {
		case *compiler.LabelEntry:
			switch e.(*compiler.LabelEntry).Name {
			case string(compiler.StringLabel), string(compiler.MethodNameLabel), string(compiler.ObjectNameLabel):
				result.addOperation(&PushStringOperation{index: e.(*compiler.LabelEntry).Value})
			default:
				labelPos[e.(*compiler.LabelEntry).Name] = len(result.operations) - 1
			}
		case *compiler.PopToRegisterEntry:
			result.addOperation(&PopToRegisterOperation{registerType: StringRegisterType, index: e.(*compiler.PopToRegisterEntry).Register})
		case *compiler.PushContextEntry:
			result.addOperation(&PushContextOperation{contextNameIndex: e.(*compiler.PushContextEntry).Name})
		case *compiler.MethodCallEntry:
			result.addOperation(&MethodCallOperation{argumentCount: 1})
		case *compiler.OperationEntry:
			switch e.(*compiler.OperationEntry).Operation {
			case compiler.OperationAdd:
				result.addOperation(&AddOperation{})
			case compiler.OperationCompare:
				result.addOperation(&EqualOperation{})
			}
		case *compiler.JumpIfFalseEntry:
			posRequestingLabel[len(result.operations)] = e.(*compiler.JumpIfFalseEntry).GetLabel()
			result.addOperation(&JumpIfFalseOperation{})
		case *compiler.JumpEntry:
			posRequestingLabel[len(result.operations)] = e.(*compiler.JumpEntry).GetLabel()
			result.addOperation(&JumpOperation{})
		case *compiler.PushFromRegisterEntry:
			result.addOperation(&PushFromRegisterOperation{registerType: StringRegisterType, index: e.(*compiler.PushFromRegisterEntry).Register})
		}
	}

	for pos, label := range posRequestingLabel {
		switch result.operations[pos].(type) {
		case *JumpIfFalseOperation:
			result.operations[pos].(*JumpIfFalseOperation).target = labelPos[label]
		case *JumpOperation:
			result.operations[pos].(*JumpOperation).target = labelPos[label]
		}

	}

	return result
}

func (m *vmMethod) GetStrings() []string {
	return m.strings
}

func (m *vmMethod) addOperation(p Operation) {
	m.operations = append(m.operations, p)
}
