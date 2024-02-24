package vm

import "log"

// div(a, b)   | StringValue      		     | ObjectValue                 | BooleanValue                    | NumberValue
// StringValue | unsupportedDivision(a,b)     | unsupportedDivision(a,b)    | unsupportedDivision(a,b)        | unsupportedDivision(a,b)
// ObjectValue | unsupportedDivision(a,b)     | unsupportedDivision(a,b)    | unsupportedDivision(a,b)        | unsupportedDivision(a,b)
// BooleanValue| unsupportedDivision(a,b)     | unsupportedDivision(a,b)    | unsupportedDivision(a,b)        | unsupportedDivision(a,b)
// NumberValue | unsupportedDivision(a,b)     | unsupportedDivision(a,b)    | unsupportedDivision(a,b)        | a / b

func divide(a Value, b Value) Value {
	switch a.(type) {
	case StringValue:
		return unsupportedDivision(a, b)
	case ObjectValue:
		return unsupportedDivision(a, b)
	case BooleanValue:
		return unsupportedDivision(a, b)
	case NumberValue:
		if b, ok := b.(NumberValue); ok {
			return NewNumberValue(a.(NumberValue).Value / b.Value)
		}
		return unsupportedDivision(a, b)
	}
	return nil
}

func unsupportedDivision(a Value, b Value) Value {
	log.Panicln("Division not supported between", a, "and", b)
	return nil
}
