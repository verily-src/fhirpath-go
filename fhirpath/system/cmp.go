package system

import "reflect"

// Equal compares two FHIRPath System types for equality. This uses standard
// equality semantics and will return true if the value should yield a value
// that is true, and false otherwise.
//
// This is effectively sugar over calling:
//
//	result, ok := TryEqual(lhs, rhs)
//	return result && ok
//
// See https://hl7.org/fhirpath/n1/#equality
func Equal(lhs, rhs Any) bool {
	result, ok := TryEqual(lhs, rhs)
	return ok && result
}

// TryEqual compares two FHIRPath System types for equality. This returns a
// value if and only if the comparison of the underlying System types should
// also yield a value, as defined in FHIRPath's equality operation.
//
// See https://hl7.org/fhirpath/n1/#equality
//
// For system types that define a custom "Equal" function, this will call the
// underlying function. For system types that define a custom "TryEqual" function
// this will call the underlying function. Otherwise, this will compare the raw
// representation instead.
func TryEqual(lhs, rhs Any) (bool, bool) {
	lhs, rhs = Normalize(lhs, rhs), Normalize(rhs, lhs)
	if result, has, ok := callTryEqual(lhs, rhs); ok {
		return result, has
	}
	if result, ok := callEqual(lhs, rhs); ok {
		return result, true
	}
	return lhs == rhs, true
}

func callTryEqual(lhs, rhs Any) (bool, bool, bool) {
	if eq, ok := reflect.TypeOf(lhs).MethodByName("TryEqual"); ok {
		funcType := eq.Func.Type()
		arg0, arg1 := funcType.In(0), funcType.In(1)
		if arg0 != reflect.TypeOf(lhs) || arg1.ConvertibleTo(reflect.TypeOf(rhs)) {
			return false, false, true
		}
		args := []reflect.Value{
			reflect.ValueOf(lhs),
			reflect.ValueOf(rhs),
		}
		result := eq.Func.Call(args)
		return result[0].Bool(), result[1].Bool(), true
	}
	return false, false, false
}

func callEqual(lhs, rhs Any) (bool, bool) {
	if got, ok := callBinaryComparator("Equal", lhs, rhs); ok {
		return got.(bool), true
	}
	return false, false
}

// callBinaryComparator invokes a binary comparator operator
func callBinaryComparator(name string, lhs, rhs any) (any, bool) {
	if eq, ok := reflect.TypeOf(lhs).MethodByName(name); ok {
		args := []reflect.Value{
			reflect.ValueOf(lhs),
			reflect.ValueOf(rhs),
		}
		funcType := eq.Func.Type()
		arg0, arg1 := funcType.In(0), funcType.In(1)
		if arg0 != reflect.TypeOf(lhs) || arg1.ConvertibleTo(reflect.TypeOf(rhs)) {
			return nil, true
		}
		result := eq.Func.Call(args)
		return result[0].Bool(), true
	}
	return nil, false
}
