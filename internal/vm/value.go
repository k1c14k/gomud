package vm

import "log"

type Value interface {
	AddValue(v Value) Value
	String() string
	isTruthy() bool
	equalValue(b Value) Value
}

type StringValue struct {
	Value string
}

func (s StringValue) AddValue(v Value) Value {
	return NewStringValue(s.Value + v.String())
}

func (s StringValue) String() string {
	return s.Value
}

func (s StringValue) isTruthy() bool {
	return s.Value != ""
}

func (s StringValue) equalValue(b Value) Value {
	if sv, ok := b.(*StringValue); ok {
		return BooleanValue{Value: s.Value == sv.Value}
	}
	return BooleanValue{Value: false}
}

func NewStringValue(value string) *StringValue {
	return &StringValue{Value: value}
}

type ObjectValue struct {
	value *Object
}

func (o ObjectValue) isTruthy() bool {
	return o.value != nil
}

func (o ObjectValue) equalValue(b Value) Value {
	if ov, ok := b.(*ObjectValue); ok {
		return BooleanValue{Value: o.value == ov.value}
	}
	return BooleanValue{Value: false}
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

type BooleanValue struct {
	Value bool
}

func (b BooleanValue) AddValue(v Value) Value {
	return BooleanValue{b.Value || v.isTruthy()}
}

func (b BooleanValue) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

func (b BooleanValue) isTruthy() bool {
	return b.Value
}

func (b BooleanValue) equalValue(v Value) Value {
	if bv, ok := v.(BooleanValue); ok {
		return BooleanValue{Value: b.Value == bv.Value}
	}
	return BooleanValue{Value: false}
}
