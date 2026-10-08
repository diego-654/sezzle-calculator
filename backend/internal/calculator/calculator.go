// Package calculator implements the arithmetic operations of the service.
// It is pure domain logic: no HTTP, no I/O, and every function returns
// (result, error) instead of panicking.
package calculator

import "math"

// Add returns a + b.
func Add(a, b float64) (float64, error) {
	if err := validateOperands(a, b); err != nil {
		return 0, err
	}
	return checkResult(a + b)
}

// Subtract returns a - b.
func Subtract(a, b float64) (float64, error) {
	if err := validateOperands(a, b); err != nil {
		return 0, err
	}
	return checkResult(a - b)
}

// Multiply returns a * b.
func Multiply(a, b float64) (float64, error) {
	if err := validateOperands(a, b); err != nil {
		return 0, err
	}
	return checkResult(a * b)
}

// Divide returns a / b. It returns ErrDivisionByZero when b is zero.
func Divide(a, b float64) (float64, error) {
	if err := validateOperands(a, b); err != nil {
		return 0, err
	}
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return checkResult(a / b)
}

// Power returns a raised to the power b.
// Zero raised to a negative exponent is a division by zero, and a negative
// base with a fractional exponent has no real result.
func Power(a, b float64) (float64, error) {
	if err := validateOperands(a, b); err != nil {
		return 0, err
	}
	if a == 0 && b < 0 {
		return 0, ErrDivisionByZero
	}
	if a < 0 && b != math.Trunc(b) {
		return 0, ErrInvalidExponent
	}
	return checkResult(math.Pow(a, b))
}

// Sqrt returns the square root of a. It returns ErrNegativeSquareRoot
// when a is negative.
func Sqrt(a float64) (float64, error) {
	if err := validateOperands(a); err != nil {
		return 0, err
	}
	if a < 0 {
		return 0, ErrNegativeSquareRoot
	}
	return checkResult(math.Sqrt(a))
}

// Percent returns b percent of a, that is a * b / 100.
// For example, Percent(200, 15) returns 30.
func Percent(a, b float64) (float64, error) {
	if err := validateOperands(a, b); err != nil {
		return 0, err
	}
	return checkResult(a * b / 100)
}

// validateOperands rejects NaN and ±Inf inputs, which are not valid numbers
// for a calculator.
func validateOperands(nums ...float64) error {
	for _, n := range nums {
		if math.IsNaN(n) || math.IsInf(n, 0) {
			return ErrInvalidOperand
		}
	}
	return nil
}

// checkResult rejects non-finite results and normalizes -0 to 0 so the API
// never returns "-0".
func checkResult(r float64) (float64, error) {
	if math.IsInf(r, 0) {
		return 0, ErrOverflow
	}
	if math.IsNaN(r) {
		return 0, ErrUndefinedResult
	}
	if r == 0 {
		return 0, nil
	}
	return r, nil
}
