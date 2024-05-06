package impl

import "errors"

// Error constants
var (
	ErrWrongArity        = errors.New("incorrect function arity")
	ErrInvalidReturnType = errors.New("invalid return type")
)
