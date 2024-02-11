package compiler

import (
	"goMud/internal/gmsl"
	"strconv"
)

type AssemblyEntry interface {
	String() string
	GetSource() gmsl.Token
}

type LabelEntry struct {
	AssemblyEntry
	Name   string
	Value  int
	source gmsl.Token
}

type PopToRegisterEntry struct {
	AssemblyEntry
	Register RegisterReference
	source   gmsl.Token
}

type PushContextEntry struct {
	AssemblyEntry
	Name   int
	source gmsl.Token
}

type MethodCallEntry struct {
	AssemblyEntry
	source gmsl.Token
}

type OperationType int

const (
	OperationAdd OperationType = iota
)

type OperationEntry struct {
	AssemblyEntry
	Operation OperationType
	source    gmsl.Token
}

type ReturnEntry struct {
	AssemblyEntry
	source gmsl.Token
}

func NewReturnEntry(source gmsl.Token) AssemblyEntry {
	return &ReturnEntry{source: source}
}

func (r *ReturnEntry) String() string {
	return "RETURN"
}

func (r *ReturnEntry) GetSource() gmsl.Token {
	return r.source
}

var tokenValueToOperationType = map[string]OperationType{
	"+": OperationAdd,
}

func NewOperationEntry(token gmsl.Token) AssemblyEntry {
	operation := tokenValueToOperationType[token.Value]
	return &OperationEntry{Operation: operation, source: token}
}

var operationTypeToString = map[OperationType]string{
	OperationAdd: "ADD",
}

func (o *OperationEntry) String() string {
	return operationTypeToString[o.Operation]
}

func NewMethodCallEntry(token gmsl.Token) AssemblyEntry {
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

func (p *PushContextEntry) GetSource() gmsl.Token {
	return p.source
}

func NewPushContextEntry(name int, source gmsl.Token) AssemblyEntry {
	return &PushContextEntry{Name: name, source: source}
}

func NewLabelEntry(label string, reference int, source gmsl.Token) AssemblyEntry {
	return &LabelEntry{Name: label, Value: reference, source: source}
}

func NewPopToRegisterEntry(r RegisterReference, token gmsl.Token) AssemblyEntry {
	return &PopToRegisterEntry{Register: r, source: token}
}
