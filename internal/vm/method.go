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
		if e.GetLabel() != nil {
			labelPos[*e.GetLabel()] = len(result.operations)
		}
		switch e.GetOpCode() {
		case compiler.OpReturn:
			result.addOperation(&ReturnOperation{})
		case compiler.OpAdd:
			result.addOperation(&AddOperation{})
		case compiler.OpCmp:
			result.addOperation(&EqualOperation{})
		case compiler.OpCall:
			result.addOperation(&MethodCallOperation{argumentCount: 1})
		case compiler.OpPushContext:
			result.addOperation(&PushContextOperation{contextNameIndex: e.GetArgument()})
		case compiler.OpPopToRegister:
			result.addOperation(&PopToRegisterOperation{registerType: StringRegisterType, index: e.GetArgument()})
		case compiler.OpJumpIfFalse:
			posRequestingLabel[len(result.operations)] = *e.GetTargetLabel()
			result.addOperation(&JumpIfFalseOperation{})
		case compiler.OpJump:
			posRequestingLabel[len(result.operations)] = *e.GetTargetLabel()
			result.addOperation(&JumpOperation{})
		case compiler.OpPushFromRegister:
			result.addOperation(&PushFromRegisterOperation{registerType: StringRegisterType, index: e.GetArgument()})
		case compiler.OpPushString:
			result.addOperation(&PushStringOperation{index: e.GetArgument()})
		case compiler.OpNoOp:
			// Do nothing
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
