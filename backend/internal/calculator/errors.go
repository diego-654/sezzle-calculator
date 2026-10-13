package calculator

import "errors"

var (
	// General errors
	ErrInvalidOperation = errors.New("invalid operation")
	ErrInvalidInput     = errors.New("invalid input")
	ErrInvalidOperand   = errors.New("invalid operand")

	// Arithmetic errors
	ErrDivisionByZero     = errors.New("division by zero")
	ErrNegativeSquareRoot = errors.New("cannot take square root of a negative number")
	ErrInvalidExponent    = errors.New("invalid exponent")

	// Result errors
	ErrUndefinedResult = errors.New("undefined result")
	ErrOverflow        = errors.New("result overflow")
)
