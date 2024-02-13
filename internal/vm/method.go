package vm

import (
	"goMud/internal/gmsl/compiler"
	"log"
)

type Method interface {
	GetArgumentCount() int
	GetReturnValueCount() int
}

type vmMethod struct {
	argumentCount    int
	returnValueCount int
	operations       []Operation
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
	m := make(map[string]Method)

	start := 0

	for start < len(aOut.Entries) {
		name, method, nextStart := NewMethodFromAssembly(aOut, start, aOut.MethodArgumentCounts)
		m[name] = method
		start = nextStart
	}

	return m
}

func NewMethodFromAssembly(aOut *compiler.Assembly, start int, counts map[string]int) (string, Method, int) {
	pos := start

	var name string
	result := &vmMethod{}
	labelPos := make(map[string]int)
	posRequestingLabel := make(map[int]string)

	for pos < len(aOut.Entries) {
		switch e := aOut.Entries[pos].(type) {
		case *compiler.LabelEntry:
			name = processLabelEntry(aOut, *e, name, result, labelPos)
		case *compiler.PopToRegisterEntry:
			processPopToRegisterEntry(*e, result)
		case *compiler.PushContextEntry:
			result.operations = append(result.operations, &PushContextOperation{contextNameIndex: e.Name})
		case *compiler.MethodCallEntry:
			result.operations = append(result.operations, &MethodCallOperation{argumentCount: 1})
		case *compiler.OperationEntry:
			processOperation(*e, result)
		case *compiler.JumpIfFalseEntry:
			result.operations = append(result.operations, &JumpIfFalseOperation{})
			posRequestingLabel[len(result.operations)-1] = e.GetLabel()
		case *compiler.JumpEntry:
			result.operations = append(result.operations, &JumpOperation{})
			posRequestingLabel[len(result.operations)-1] = e.GetLabel()
		case *compiler.PushFromRegisterEntry:
			processPushFromRegister(*e, result)
		}
		pos++
	}

	for pos, label := range posRequestingLabel {
		if target, ok := labelPos[label]; ok {
			switch o := result.operations[pos].(type) {
			case *JumpIfFalseOperation:
				o.target = target
			case *JumpOperation:
				o.target = target
			}
		} else {
			log.Panicln("Unknown label", label)
		}
	}

	result.argumentCount = counts[name]

	return name, result, pos + 1
}

func processPushFromRegister(e compiler.PushFromRegisterEntry, result *vmMethod) {
	var registerType RegisterType
	switch e.Register.Typ {
	case compiler.StringRegister:
		registerType = StringRegisterType
	}
	result.operations = append(result.operations, &PushFromRegisterOperation{registerType: registerType, index: e.Register.Index})
}

func processOperation(e compiler.OperationEntry, result *vmMethod) {
	switch e.Operation {
	case compiler.OperationAdd:
		result.operations = append(result.operations, &AddOperation{})
	case compiler.OperationCompare:
		result.operations = append(result.operations, &EqualOperation{})
	}
}

func processPopToRegisterEntry(e compiler.PopToRegisterEntry, result *vmMethod) {
	var registerType RegisterType
	switch e.Register.Typ {
	case compiler.StringRegister:
		registerType = StringRegisterType
	default:
		log.Panicln("Unknown register type", e.Register.Typ)
	}
	result.operations = append(result.operations, &PopToRegisterOperation{registerType: registerType, index: e.Register.Index})
}

func processLabelEntry(aOut *compiler.Assembly, e compiler.LabelEntry, name string, result *vmMethod, labelPos map[string]int) string {
	if e.Name == ".function_name" {
		name = aOut.Consts[e.Value]
	} else if e.Name == ".string" || e.Name == ".method_name" {
		result.operations = append(result.operations, &PushStringOperation{index: e.Value})
	} else {
		labelPos[e.Name] = len(result.operations) - 1
	}
	return name
}

func (m *vmMethod) Execute(_ *ExecutionFrame) {
	log.Panicln("Method not implemented")
}
