package compiler

import (
	"goMud/internal/gmsl/lexer"
	"strconv"
)

type AssemblyEntry interface {
	String() string
	GetSource() lexer.Token
}

type LabelEntry struct {
	AssemblyEntry
	Name   string
	Value  int
	source lexer.Token
}

type PopToRegisterEntry struct {
	AssemblyEntry
	Register int
	source   lexer.Token
}

type PushContextEntry struct {
	AssemblyEntry
	Name   int
	source lexer.Token
}

type MethodCallEntry struct {
	AssemblyEntry
	source lexer.Token
}

type OperationType int

const (
	OperationAdd OperationType = iota
	OperationCompare
)

type OperationEntry struct {
	AssemblyEntry
	Operation OperationType
	source    lexer.Token
}

type ReturnEntry struct {
	AssemblyEntry
	source lexer.Token
}

type JumpIfFalseEntry struct {
	source lexer.Token
	label  string
}

type JumpEntry struct {
	source lexer.Token
	label  string
}

type PushFromRegisterEntry struct {
	Register int
	source   lexer.Token
}

func NewReturnEntry(source lexer.Token) AssemblyEntry {
	return &ReturnEntry{source: source}
}

func (r *ReturnEntry) String() string {
	return "RETURN"
}

func (r *ReturnEntry) GetSource() lexer.Token {
	return r.source
}

var tokenValueToOperationType = map[string]OperationType{
	"+":  OperationAdd,
	"==": OperationCompare,
}

func NewOperationEntry(token lexer.Token) AssemblyEntry {
	operation := tokenValueToOperationType[token.GetRawValue()]
	return &OperationEntry{Operation: operation, source: token}
}

var operationTypeToString = map[OperationType]string{
	OperationAdd:     "ADD",
	OperationCompare: "CMP",
}

func (o *OperationEntry) String() string {
	return operationTypeToString[o.Operation]
}

func NewMethodCallEntry(token lexer.Token) AssemblyEntry {
	return &MethodCallEntry{source: token}
}

func (m *MethodCallEntry) String() string {
	return "MCALL"
}

func (l *LabelEntry) String() string {
	return l.Name + " $" + strconv.Itoa(l.Value)
}

func (p *PopToRegisterEntry) String() string {
	return "RPOP " + strconv.Itoa(p.Register)
}

func (p *PushContextEntry) String() string {
	return "CPUSH $" + strconv.Itoa(p.Name)
}

func (p *PushContextEntry) GetSource() lexer.Token {
	return p.source
}

func NewPushContextEntry(name int, source lexer.Token) AssemblyEntry {
	return &PushContextEntry{Name: name, source: source}
}

func NewLabelEntry(label string, reference int, source lexer.Token) AssemblyEntry {
	return &LabelEntry{Name: label, Value: reference, source: source}
}

func NewPopToRegisterEntry(r int, token lexer.Token) AssemblyEntry {
	return &PopToRegisterEntry{Register: r, source: token}
}

func NewJumpIfFalseEntry(label string, source lexer.Token) AssemblyEntry {
	return &JumpIfFalseEntry{label: label, source: source}
}

func NewJumpEntry(label string, source lexer.Token) AssemblyEntry {
	return &JumpEntry{label: label, source: source}
}

func (j *JumpIfFalseEntry) String() string {
	return "JMPF " + j.label
}

func (j *JumpEntry) String() string {
	return "JMP " + j.label
}

func (j *JumpIfFalseEntry) GetSource() lexer.Token {
	return j.source
}

func (j *JumpEntry) GetSource() lexer.Token {
	return j.source
}

func (j *JumpIfFalseEntry) GetLabel() string {
	return j.label
}

func (j *JumpEntry) GetLabel() string {
	return j.label
}

func (p *PushFromRegisterEntry) String() string {
	return "RPUSH " + strconv.Itoa(p.Register)
}

func NewPushFromRegisterEntry(register int, source lexer.Token) AssemblyEntry {
	return &PushFromRegisterEntry{Register: register, source: source}
}

func (p *PushFromRegisterEntry) GetSource() lexer.Token {
	return p.source
}
