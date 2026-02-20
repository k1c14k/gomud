package vm

import "github.com/sirupsen/logrus"

// sub(a, b)   | StringValue      		     | ObjectValue                 | BooleanValue                    | NumberValue
// StringValue | unsupportedSubtraction(a,b) | unsupportedSubtraction(a,b) | unsupportedSubtraction(a,b)     | unsupportedSubtraction(a,b)
// ObjectValue | unsupportedSubtraction(a,b) | unsupportedSubtraction(a,b) | unsupportedSubtraction(a,b)     | unsupportedSubtraction(a,b)
// BooleanValue| unsupportedSubtraction(a,b) | unsupportedSubtraction(a,b) | unsupportedSubtraction(a,b)     | unsupportedSubtraction(a,b)
// NumberValue | unsupportedSubtraction(a,b) | unsupportedSubtraction(a,b) | unsupportedSubtraction(a,b)     | a - b

func subtract(a Value, b Value) Value {
	switch a.(type) {
	case StringValue:
		return unsupportedSubtraction(a, b)
	case ObjectValue:
		return unsupportedSubtraction(a, b)
	case BooleanValue:
		return unsupportedSubtraction(a, b)
	case NumberValue:
		if b, ok := b.(NumberValue); ok {
			return NewNumberValue(a.(NumberValue).Value - b.Value)
		}
		return unsupportedSubtraction(a, b)
	}
	return nil
}

func unsupportedSubtraction(a Value, b Value) Value {
	logrus.Panic("Subtraction not supported between", a, "and", b)
	return nil
}
