package vm

import (
	"goMud/internal/gmsl/compiler"
	"log"
)

type Method interface {
	Execute(ef *ExecutionFrame)
}

type vmMethod struct {
	operations []Operation
}

func NewMethodsFromAssembly(aOut *compiler.Assembly) map[string]Method {
	m := make(map[string]Method)

	start := 0

	for start < len(aOut.Entries) {
		name, method, nextStart := NewMethodFromAssembly(aOut, start)
		m[name] = method
		start = nextStart
	}

	return m
}

func NewMethodFromAssembly(aOut *compiler.Assembly, start int) (string, Method, int) {
	pos := start

	var name string
	result := &vmMethod{}

	for pos < len(aOut.Entries) {
		switch e := aOut.Entries[pos].(type) {
		case *compiler.LabelEntry:
			if e.Name == ".function_name" {
				name = aOut.Consts[e.Value]
			} else if e.Name == ".string" || e.Name == ".method_name" {
				result.operations = append(result.operations, &PushStringOperation{index: e.Value})
			}
		case *compiler.PopToRegisterEntry:
			var registerType RegisterType
			switch e.Register.Typ {
			case compiler.StringRegister:
				registerType = StringRegisterType
			default:
				log.Panicln("Unknown register type", e.Register.Typ)
			}
			result.operations = append(result.operations, &PopToRegisterOperation{registerType: registerType, index: e.Register.Index})
		case *compiler.PushContextEntry:
			result.operations = append(result.operations, &PushContextOperation{contextNameIndex: e.Name})
		case *compiler.MethodCallEntry:
			result.operations = append(result.operations, &MethodCallOperation{argumentCount: 1})
		case *compiler.OperationEntry:
			switch e.Operation {
			case compiler.OperationAdd:
				result.operations = append(result.operations, &AddOperation{})
			}
		}
		pos++
	}

	return name, result, pos + 1
}

func (m *vmMethod) Execute(_ *ExecutionFrame) {
	log.Panicln("Method not implemented")
}
