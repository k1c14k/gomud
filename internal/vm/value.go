package vm

import "log"

type Value interface {
	AddValue(v Value) Value
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

func NewStringValue(value string) *StringValue {
	return &StringValue{Value: value}
}

type ObjectValue struct {
	value *Object
}

func (o ObjectValue) AddValue(_ Value) Value {
	log.Panicln("Cannot add values")
	return nil
}
