package compiler

import (
	"bytes"
	"goMud/internal/gmsl/lexer"
	"strconv"
)

type OpCode int

const (
	OpReturn OpCode = iota
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod
	OpCmp
	OpCall
	OpPushContext
	OpPopToRegister
	OpJumpIfFalse
	OpJump
	OpPushFromRegister
	OpPushString
	OpPushNumber
	OpNoOp
)

var opCodeString = map[OpCode]string{
	OpReturn:           "RET",
	OpAdd:              "ADD",
	OpSub:              "SUB",
	OpMul:              "MUL",
	OpDiv:              "DIV",
	OpMod:              "MOD",
	OpCmp:              "CMP",
	OpCall:             "CALL",
	OpPushContext:      "PUCX",
	OpPopToRegister:    "POPR",
	OpJumpIfFalse:      "JMPF",
	OpJump:             "JMP",
	OpPushFromRegister: "PURE",
	OpPushString:       "PUSC",
	OpPushNumber:       "PUSN",
	OpNoOp:             "NOOP",
}

func (o OpCode) String() string {
	return opCodeString[o]
}

type AssemblyEntry struct {
	label         *string
	opCode        OpCode
	argument      *int
	labelArgument *string
	source        lexer.Token
}

func (a *AssemblyEntry) String() string {
	var buf bytes.Buffer
	if a.label != nil {
		buf.WriteString(*a.label)
		buf.WriteString(":\n")
	}
	buf.WriteString(a.opCode.String())
	if a.argument != nil {
		buf.WriteString(" ")
		buf.WriteString(strconv.Itoa(*a.argument))
	}
	if a.labelArgument != nil {
		buf.WriteString(" ")
		buf.WriteString(*a.labelArgument)
	}
	return buf.String()
}

func (a *AssemblyEntry) GetSource() lexer.Token {
	return a.source
}

func (a *AssemblyEntry) GetOpCode() OpCode {
	return a.opCode
}

func (a *AssemblyEntry) GetLabel() *string {
	return a.label
}

func (a *AssemblyEntry) GetArgument() int {
	return *a.argument
}

func (a *AssemblyEntry) GetTargetLabel() *string {
	return a.labelArgument
}

func NewReturnEntry(label *string, source lexer.Token) *AssemblyEntry {
	return &AssemblyEntry{label: label, opCode: OpReturn, source: source}
}

var tokenToOpCode = map[string]OpCode{
	"+":  OpAdd,
	"-":  OpSub,
	"*":  OpMul,
	"/":  OpDiv,
	"%":  OpMod,
	"==": OpCmp,
}

func NewOperationEntry(label *string, source lexer.Token) *AssemblyEntry {
	operation := tokenToOpCode[source.GetRawValue()]
	return &AssemblyEntry{
		label:    label,
		opCode:   operation,
		argument: nil,
		source:   source,
	}
}

func NewCallEntry(label *string, token lexer.Token) *AssemblyEntry {
	return &AssemblyEntry{label: label, opCode: OpCall, argument: nil, source: token}
}

func NewPushContextEntry(label *string, nameIdx int, source lexer.Token) *AssemblyEntry {
	return &AssemblyEntry{label: label, opCode: OpPushContext, argument: &nameIdx, source: source}
}

func NewPopToRegisterEntry(label *string, idx int, token lexer.Token) *AssemblyEntry {
	return &AssemblyEntry{label: label, opCode: OpPopToRegister, argument: &idx, source: token}
}

func NewJumpIfFalseEntry(label *string, target string, source lexer.Token) *AssemblyEntry {
	return &AssemblyEntry{label: label, opCode: OpJumpIfFalse, labelArgument: &target, source: source}
}

func NewJumpEntry(label *string, target string, source lexer.Token) *AssemblyEntry {
	return &AssemblyEntry{label: label, opCode: OpJump, labelArgument: &target, source: source}
}

func NewPushFromRegisterEntry(label *string, register int, source lexer.Token) *AssemblyEntry {
	return &AssemblyEntry{label: label, opCode: OpPushFromRegister, argument: &register, source: source}
}

func NewPushStringEntry(label *string, stringIdx int, source lexer.Token) *AssemblyEntry {
	return &AssemblyEntry{label: label, opCode: OpPushString, argument: &stringIdx, source: source}
}

func NewNoOpEntry(label *string, source lexer.Token) *AssemblyEntry {
	return &AssemblyEntry{label: label, opCode: OpNoOp, source: source}
}

func NewPushNumberEntry(label *string, value int, source lexer.Token) *AssemblyEntry {
	return &AssemblyEntry{label: label, opCode: OpPushNumber, argument: &value, source: source}
}
