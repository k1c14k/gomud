package vm

import "log"

type Value interface {
	AddValue(v Value) Value
	String() string
}

type StringValue struct {
	Value string
}

func (s StringValue) AddValue(v Value) Value {
	if sv, ok := v.(*StringValue); ok {
		return &StringValue{Value: s.Value + sv.Value}
	}
	log.Panicln("Cannot add values")
	return nil
}

func (s StringValue) String() string {
	return s.Value
}

func NewStringValue(value string) *StringValue {
	return &StringValue{Value: value}
}

type ObjectValue struct {
	value *Object
}

func NewObjectValue(value *Object) *ObjectValue {
	return &ObjectValue{value: value}
}

func (o ObjectValue) AddValue(_ Value) Value {
	log.Panicln("Cannot add values")
	return nil
}

func (o ObjectValue) String() string {
	return "Object"
}
