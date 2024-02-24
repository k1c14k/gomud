package vm

import "log"

// add(a,b)    | StringValue      		   | ObjectValue              | BooleanValue              | NumberValue
// StringValue | concatenate(a,b)		   | unsupportedAddition(a,b) | concatenate(a,b.String()) | concatenate(a,b.String())
// ObjectValue | unsupportedAddition(a,b)  | unsupportedAddition(a,b) | unsupportedAddition(a,b)  | unsupportedAddition(a,b)
// BooleanValue| concatenate(a.String(),b) | unsupportedAddition(a,b) | or(a, b)                  | or(a, b.isTruthy())
// NumberValue | concatenate(a.String(),b) | unsupportedAddition(a,b) | or(a.isTruthy(), b)       | a + b

func add(a Value, b Value) Value {
	switch a.(type) {
	case StringValue:
		switch b.(type) {
		case StringValue:
			return concatenate(a, b)
		case ObjectValue:
			return unsupportedAddition(a, b)
		case BooleanValue:
			return concatenate(a, NewStringValue(b.String()))
		case NumberValue:
			return concatenate(a, NewStringValue(b.String()))
		}
	case ObjectValue:
		return unsupportedAddition(a, b)
	case BooleanValue:
		switch b.(type) {
		case StringValue:
			return concatenate(NewStringValue(a.String()), b)
		case ObjectValue:
			return unsupportedAddition(a, b)
		case BooleanValue:
			return or(a, b)
		case NumberValue:
			return or(a, BooleanValue{Value: b.isTruthy()})
		}
	case NumberValue:
		switch b.(type) {
		case StringValue:
			return concatenate(NewStringValue(a.String()), b)
		case ObjectValue:
			return unsupportedAddition(a, b)
		case BooleanValue:
			return or(BooleanValue{Value: a.isTruthy()}, b)
		case NumberValue:
			return NewNumberValue(a.(NumberValue).Value + b.(NumberValue).Value)
		}
	}
	return nil
}

func or(a Value, b Value) Value {
	return BooleanValue{Value: a.isTruthy() || b.isTruthy()}
}

func unsupportedAddition(a Value, b Value) Value {
	log.Panicln("Addition not supported between", a, "and", b)
	return nil
}

func concatenate(a Value, b Value) Value {
	return NewStringValue(a.String() + b.String())
}
