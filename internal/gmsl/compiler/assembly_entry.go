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
	Register RegisterReference
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
	"+": OperationAdd,
}

func NewOperationEntry(token lexer.Token) AssemblyEntry {
	operation := tokenValueToOperationType[token.Value]
	return &OperationEntry{Operation: operation, source: token}
}

var operationTypeToString = map[OperationType]string{
	OperationAdd: "ADD",
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

var registrySymbol = map[RegisterType]string{
	StringRegister: "S",
}

func (p *PopToRegisterEntry) String() string {
	return "RPOP " + registrySymbol[p.Register.Typ] + strconv.Itoa(p.Register.Index)
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

func NewPopToRegisterEntry(r RegisterReference, token lexer.Token) AssemblyEntry {
	return &PopToRegisterEntry{Register: r, source: token}
}
