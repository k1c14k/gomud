package vm

import (
	"strconv"
)

type Value interface {
	Add(v Value) Value
	String() string
	isTruthy() bool
	equalValue(b Value) Value
	Subtract(v Value) Value
	Multiply(v Value) Value
	Divide(v Value) Value
	Modulo(v Value) Value
}

type StringValue struct {
	Value string
}

func (s StringValue) Modulo(v Value) Value {
	return modulo(s, v)
}

func (s StringValue) Divide(v Value) Value {
	return divide(s, v)
}

func (s StringValue) Multiply(v Value) Value {
	return multiply(s, v)
}

func (s StringValue) Subtract(v Value) Value {
	return subtract(s, v)
}

func (s StringValue) Add(v Value) Value {
	return add(s, v)
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

func (o ObjectValue) Add(v Value) Value {
	return add(o, v)
}

func (o ObjectValue) Subtract(v Value) Value {
	return subtract(o, v)
}

func (o ObjectValue) Divide(v Value) Value {
	return divide(o, v)
}

func (o ObjectValue) String() string {
	return "Object"
}

func (o ObjectValue) Multiply(v Value) Value {
	return multiply(o, v)
}

func (o ObjectValue) Modulo(v Value) Value {
	return modulo(o, v)
}

type BooleanValue struct {
	Value bool
}

func (b BooleanValue) Modulo(v Value) Value {
	return modulo(b, v)
}

func (b BooleanValue) Divide(v Value) Value {
	return divide(b, v)
}

func (b BooleanValue) Multiply(v Value) Value {
	return multiply(b, v)
}

func (b BooleanValue) Subtract(v Value) Value {
	return subtract(b, v)
}

func (b BooleanValue) Add(v Value) Value {
	return add(b, v)
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

type NumberValue struct {
	Value int
}

func (n NumberValue) Modulo(v Value) Value {
	return modulo(n, v)
}

func (n NumberValue) Divide(v Value) Value {
	return divide(n, v)
}

func (n NumberValue) Multiply(v Value) Value {
	return multiply(n, v)
}

func (n NumberValue) Subtract(v Value) Value {
	return subtract(n, v)
}

func NewNumberValue(value int) NumberValue {
	return NumberValue{Value: value}
}

func (n NumberValue) String() string {
	return strconv.Itoa(n.Value)
}

func (n NumberValue) isTruthy() bool {
	return n.Value != 0
}

func (n NumberValue) equalValue(v Value) Value {
	if nv, ok := v.(NumberValue); ok {
		return BooleanValue{Value: n.Value == nv.Value}
	}
	return BooleanValue{Value: false}
}

func (n NumberValue) Add(v Value) Value {
	return add(n, v)
}
