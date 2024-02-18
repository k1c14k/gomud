package compiler

import (
	"bytes"
	"strconv"
)

type Type int

const (
	StringType Type = iota
)

type IdentifierReference struct {
	register int
	typ      Type
}

type FunctionInfo struct {
	name              string
	arguments         []Type
	returns           []Type
	strings           []string
	entries           []AssemblyEntry
	identifierNameMap map[string]IdentifierReference
}

type Assembly struct {
	functions []FunctionInfo
}

func newAssembly() *Assembly {
	return &Assembly{make([]FunctionInfo, 0)}
}

func (a *Assembly) String() string {
	var b bytes.Buffer
	for _, f := range a.functions {
		b.WriteString("Function ")
		b.WriteString(f.name)
		b.WriteString(":\n")
		b.WriteString(f.String())
	}
	return b.String()
}

func (a *Assembly) addFunction(info *FunctionInfo) {
	a.functions = append(a.functions, *info)
}

func (a *Assembly) GetFunctions() []FunctionInfo {
	return a.functions
}

func (t Type) String() string {
	switch t {
	case StringType:
		return "string"
	default:
		return "unknown"
	}
}

func newFunctionInfo(name string) *FunctionInfo {
	return &FunctionInfo{name, make([]Type, 0), make([]Type, 0), make([]string, 0), make([]AssemblyEntry, 0), make(map[string]IdentifierReference)}
}

func (f *FunctionInfo) String() string {
	var b bytes.Buffer
	b.WriteString("Arguments: ")
	for n, a := range f.arguments {
		b.WriteString(a.String())
		if n != len(f.arguments)-1 {
			b.WriteString(", ")
		}
	}
	b.WriteString("\n")
	b.WriteString("Returns: ")
	for n, r := range f.returns {
		b.WriteString(r.String())
		if n != len(f.returns)-1 {
			b.WriteString(", ")
		}
	}
	b.WriteString("Strings:\n")
	for _, s := range f.strings {
		b.WriteString("string[")
		b.WriteString(strconv.Itoa(len(s)))
		b.WriteString("]: \"")
		b.WriteString(s)
		b.WriteString("\"\n")
	}
	b.WriteString("Entries:\n")
	for _, e := range f.entries {
		b.WriteString(e.String())
		b.WriteString("\n")
	}
	return b.String()
}

func (f *FunctionInfo) addEntry(entry AssemblyEntry) {
	f.entries = append(f.entries, entry)
}

func (f *FunctionInfo) addEntries(entries []AssemblyEntry) {
	f.entries = append(f.entries, entries...)
}

func (f *FunctionInfo) addString(value string) int {
	for n, v := range f.strings {
		if v == value {
			return n
		}
	}
	f.strings = append(f.strings, value)
	return len(f.strings) - 1
}

func (f *FunctionInfo) nextEntryPost() int {
	return len(f.entries)
}

func (f *FunctionInfo) addIdentifier(value string, t Type) {
	f.identifierNameMap[value] = IdentifierReference{len(f.identifierNameMap), t}
}

func (f *FunctionInfo) getRegisterOf(value string) int {
	return f.identifierNameMap[value].register
}

func (f *FunctionInfo) GetName() string {
	return f.name
}

func (f *FunctionInfo) GetEntries() []AssemblyEntry {
	return f.entries
}

func (f *FunctionInfo) GetStrings() []string {
	return f.strings
}

func (f *FunctionInfo) GetReturnValueCount() int {
	return len(f.returns)
}

func (f *FunctionInfo) GetArgumentCount() int {
	return len(f.arguments)
}

func (f *FunctionInfo) addArgument(value string, t Type) {
	f.arguments = append(f.arguments, t)
	f.addIdentifier(value, t)
}
