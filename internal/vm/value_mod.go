package vm

import "log"

// mod(a, b)   | StringValue      		      | ObjectValue                 | BooleanValue                    | NumberValue
// StringValue | unsupportedModulo(a,b)       | unsupportedModulo(a,b)      | unsupportedModulo(a,b)          | unsupportedModulo(a,b)
// ObjectValue | unsupportedModulo(a,b)       | unsupportedModulo(a,b)      | unsupportedModulo(a,b)          | unsupportedModulo(a,b)
// BooleanValue| unsupportedModulo(a,b)       | unsupportedModulo(a,b)      | unsupportedModulo(a,b)          | unsupportedModulo(a,b)
// NumberValue | unsupportedModulo(a,b)       | unsupportedModulo(a,b)      | unsupportedModulo(a,b)          | a % b

func modulo(a Value, b Value) Value {
	switch a.(type) {
	case StringValue:
		return unsupportedModulo(a, b)
	case ObjectValue:
		return unsupportedModulo(a, b)
	case BooleanValue:
		return unsupportedModulo(a, b)
	case NumberValue:
		if b, ok := b.(NumberValue); ok {
			return NewNumberValue(a.(NumberValue).Value % b.Value)
		}
		return unsupportedModulo(a, b)
	}
	return nil
}

func unsupportedModulo(a Value, b Value) Value {
	log.Panicln("Modulo not supported between", a, "and", b)
	return nil
}
