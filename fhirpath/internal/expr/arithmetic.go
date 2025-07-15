package expr

import (
	"fmt"

	"github.com/verily-src/fhirpath-go/fhirpath/system"
)

// EvaluateAdd takes in two system types, and calls the appropriate Add method.
func EvaluateAdd(lhs, rhs system.Any) (system.Any, error) {
	switch left := lhs.(type) {
	case system.String:
		if right, ok := rhs.(system.String); ok {
			return left.Add(right), nil
		}
		return nil, typeMismatch(Add, lhs, rhs)
	case system.Integer:
		if right, ok := rhs.(system.Integer); ok {
			return left.Add(right)
		}
		return nil, typeMismatch(Add, lhs, rhs)
	case system.Decimal:
		if right, ok := rhs.(system.Decimal); ok {
			return left.Add(right), nil
		}
		return nil, typeMismatch(Add, lhs, rhs)
	case system.Time:
		if right, ok := rhs.(system.Quantity); ok {
			return left.Add(right)
		}
		return nil, typeMismatch(Add, lhs, rhs)
	case system.Date:
		if right, ok := rhs.(system.Quantity); ok {
			return left.Add(right)
		}
		return nil, typeMismatch(Add, lhs, rhs)
	case system.DateTime:
		if right, ok := rhs.(system.Quantity); ok {
			return left.Add(right)
		}
		return nil, typeMismatch(Add, lhs, rhs)
	case system.Quantity:
		if right, ok := rhs.(system.Quantity); ok {
			return left.Add(right)
		}
		return nil, typeMismatch(Add, lhs, rhs)
	default:
		return nil, typeMismatch(Add, lhs, rhs)
	}
}

// EvaluateSub takes in two system types, and calls the appropriate Sub method.
func EvaluateSub(lhs, rhs system.Any) (system.Any, error) {
	switch left := lhs.(type) {
	case system.Integer:
		if right, ok := rhs.(system.Integer); ok {
			return left.Sub(right)
		}
		return nil, typeMismatch(Sub, lhs, rhs)
	case system.Decimal:
		if right, ok := rhs.(system.Decimal); ok {
			return left.Sub(right), nil
		}
		return nil, typeMismatch(Sub, lhs, rhs)
	case system.Time:
		if right, ok := rhs.(system.Quantity); ok {
			return left.Sub(right)
		}
		return nil, typeMismatch(Sub, lhs, rhs)
	case system.Date:
		if right, ok := rhs.(system.Quantity); ok {
			return left.Sub(right)
		}
		return nil, typeMismatch(Sub, lhs, rhs)
	case system.DateTime:
		if right, ok := rhs.(system.Quantity); ok {
			return left.Sub(right)
		}
		return nil, typeMismatch(Sub, lhs, rhs)
	case system.Quantity:
		if right, ok := rhs.(system.Quantity); ok {
			return left.Sub(right)
		}
		return nil, typeMismatch(Sub, lhs, rhs)
	default:
		return nil, typeMismatch(Sub, lhs, rhs)
	}
}

// EvaluateMul takes in two system types, and calls the appropriate Mul method.
func EvaluateMul(lhs, rhs system.Any) (system.Any, error) {
	switch left := lhs.(type) {
	case system.Integer:
		if right, ok := rhs.(system.Integer); ok {
			return left.Mul(right)
		}
		if _, ok := rhs.(system.Quantity); ok {
			// TODO: Implement multiplication with Quantity.
			return nil, ErrToBeImplemented
		}
		return nil, typeMismatch(Mul, lhs, rhs)
	case system.Decimal:
		if right, ok := rhs.(system.Decimal); ok {
			return left.Mul(right), nil
		}
		if _, ok := rhs.(system.Quantity); ok {
			// TODO: Implement multiplication with Quantity.
			return nil, ErrToBeImplemented
		}
		return nil, typeMismatch(Mul, lhs, rhs)
	case system.Quantity:
		// TODO: Implement support for multiplication with Quantity.
		return nil, ErrToBeImplemented
	default:
		return nil, typeMismatch(Mul, lhs, rhs)
	}
}

// EvaluateDiv takes in two system types, and calls the appropriate Div method.
func EvaluateDiv(lhs, rhs system.Any) (system.Any, error) {
	switch left := lhs.(type) {
	case system.Integer:
		if right, ok := rhs.(system.Integer); ok {
			return left.Div(right), nil
		}
		if _, ok := rhs.(system.Quantity); ok {
			// TODO: Implement division with Quantity.
			return nil, ErrToBeImplemented
		}
		return nil, typeMismatch(Div, lhs, rhs)
	case system.Decimal:
		if right, ok := rhs.(system.Decimal); ok {
			return left.Div(right), nil
		}
		if _, ok := rhs.(system.Quantity); ok {
			// TODO: Implement division with Quantity.
			return nil, ErrToBeImplemented
		}
		return nil, typeMismatch(Div, lhs, rhs)
	case system.Quantity:
		// TODO: Implement support for division with Quantity.
		return nil, ErrToBeImplemented
	default:
		return nil, typeMismatch(Div, lhs, rhs)
	}
}

// EvaluateFloorDiv takes in two system types, and calls the appropriate FloorDiv method.
func EvaluateFloorDiv(lhs, rhs system.Any) (system.Any, error) {
	switch left := lhs.(type) {
	case system.Integer:
		if right, ok := rhs.(system.Integer); ok {
			return left.FloorDiv(right), nil
		}
		if _, ok := rhs.(system.Quantity); ok {
			// TODO: Implement floor division with Quantity.
			return nil, ErrToBeImplemented
		}
		return nil, typeMismatch(FloorDiv, lhs, rhs)
	case system.Decimal:
		if right, ok := rhs.(system.Decimal); ok {
			return left.FloorDiv(right)
		}
		if _, ok := rhs.(system.Quantity); ok {
			// TODO: Implement floor division with Quantity.
			return nil, ErrToBeImplemented
		}
		return nil, typeMismatch(FloorDiv, lhs, rhs)
	case system.Quantity:
		// TODO: Implement support for floor division with Quantity.
		return nil, ErrToBeImplemented
	default:
		return nil, typeMismatch(FloorDiv, lhs, rhs)
	}
}

// EvaluateMod takes in two system types, and calls the appropriate Mod method.
func EvaluateMod(lhs, rhs system.Any) (system.Any, error) {
	switch left := lhs.(type) {
	case system.Integer:
		if right, ok := rhs.(system.Integer); ok {
			return left.Mod(right), nil
		}
		if _, ok := rhs.(system.Quantity); ok {
			// TODO: Implement modulus with Quantity.
			return nil, ErrToBeImplemented
		}
		return nil, typeMismatch(Mod, lhs, rhs)
	case system.Decimal:
		if right, ok := rhs.(system.Decimal); ok {
			return left.Mod(right), nil
		}
		if _, ok := rhs.(system.Quantity); ok {
			// TODO: Implement modulus with Quantity.
			return nil, ErrToBeImplemented
		}
		return nil, typeMismatch(Mod, lhs, rhs)
	case system.Quantity:
		// TODO: Implement support for modulus with Quantity.
		return nil, ErrToBeImplemented
	default:
		return nil, typeMismatch(Mod, lhs, rhs)
	}
}

// typeMismatch generates an unsupported operation error.
func typeMismatch(op Operator, lhs, rhs system.Any) error {
	return fmt.Errorf("%w: %T %s %T", system.ErrTypeMismatch, lhs, op, rhs)
}
