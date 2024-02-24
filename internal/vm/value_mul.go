package vm

import (
	"log"
	"strings"
)

// mul(a, b)   | StringValue      		        | ObjectValue                    | BooleanValue | NumberValue
// StringValue | unsupportedMultiplication(a,b) | unsupportedMultiplication(a,b) | b ? a : ""   | repeat(a,b)
// ObjectValue | unsupportedMultiplication(a,b) | unsupportedMultiplication(a,b) | b ? a : nil  | unsupportedMultiplication(a,b)
// BooleanValue| a ? b : ""                     | a ? b : nil                    | a && b       | a ? b : 0
// NumberValue | repeat(b, a)                   | unsupportedMultiplication(a,b) | b ? a : 0    | a * b

func multiply(a Value, b Value) Value {
	switch a.(type) {
	case StringValue:
		switch b.(type) {
		case BooleanValue:
			if b.isTruthy() {
				return a
			}
			return NewStringValue("")
		case NumberValue:
			return repeat(a, b)
		}
		return unsupportedMultiplication(a, b)
	case ObjectValue:
		if b, ok := b.(BooleanValue); ok {
			if b.isTruthy() {
				return a
			}
			return NewObjectValue(nil)
		}
		return unsupportedMultiplication(a, b)
	case BooleanValue:
		value, done := multiplyBoolean(a.(BooleanValue), b)
		if done {
			return value
		}
	case NumberValue:
		value, done := multiplyNumber(a.(NumberValue), b)
		if done {
			return value
		}
	}
	return nil
}

func multiplyBoolean(a BooleanValue, b Value) (Value, bool) {
	switch b.(type) {
	case StringValue:
		if a.isTruthy() {
			return b, true
		}
		return NewStringValue(""), true
	case ObjectValue:
		if a.isTruthy() {
			return b, true
		}
		return NewObjectValue(nil), true
	case BooleanValue:
		return and(a, b), true
	case NumberValue:
		if a.isTruthy() {
			return b, true
		}
		return NewNumberValue(0), true
	}
	return nil, false
}

func multiplyNumber(a NumberValue, b Value) (Value, bool) {
	switch b.(type) {
	case StringValue:
		return repeat(b, a), true
	case ObjectValue:
		return unsupportedMultiplication(a, b), true
	case BooleanValue:
		if b.isTruthy() {
			return a, true
		}
		return NewNumberValue(0), true
	case NumberValue:
		return NewNumberValue(a.Value * b.(NumberValue).Value), true
	}
	return nil, false
}

func and(a Value, b Value) Value {
	return BooleanValue{Value: a.isTruthy() && b.isTruthy()}
}

func unsupportedMultiplication(a Value, b Value) Value {
	log.Panicln("Multiplication not supported between", a, "and", b)
	return nil
}

func repeat(a Value, b Value) Value {
	if b, ok := b.(NumberValue); ok {
		if b.Value < 0 {
			return NewStringValue("")
		}
		return NewStringValue(strings.Repeat(a.String(), b.Value))
	}
	return unsupportedMultiplication(a, b)
}
