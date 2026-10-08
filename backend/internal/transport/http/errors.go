package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/diego-654/sezzle-calculator/backend/internal/calculator"
)

// apiError 
type apiError struct {
	status int    // HTTP status code
	code   string // stable machine-readable code
}

// errorMappings lists every known domain error and its HTTP form.

var errorMappings = []struct {
	target error
	api    apiError
}{
	// Bad requests
	{calculator.ErrInvalidOperation, apiError{http.StatusNotFound, "INVALID_OPERATION"}},
	{calculator.ErrInvalidInput, apiError{http.StatusBadRequest, "INVALID_INPUT"}},
	{calculator.ErrInvalidOperand, apiError{http.StatusBadRequest, "INVALID_OPERAND"}},

	// Well-formed requests the math cannot answer.
	{calculator.ErrDivisionByZero, apiError{http.StatusUnprocessableEntity, "DIVISION_BY_ZERO"}},
	{calculator.ErrNegativeSquareRoot, apiError{http.StatusUnprocessableEntity, "NEGATIVE_SQUARE_ROOT"}},
	{calculator.ErrInvalidExponent, apiError{http.StatusUnprocessableEntity, "INVALID_EXPONENT"}},
	{calculator.ErrInvalidPercentage, apiError{http.StatusUnprocessableEntity, "INVALID_PERCENTAGE"}},
	{calculator.ErrOverflow, apiError{http.StatusUnprocessableEntity, "OVERFLOW"}},
	{calculator.ErrResultOutOfRange, apiError{http.StatusUnprocessableEntity, "RESULT_OUT_OF_RANGE"}},
	{calculator.ErrUndefinedResult, apiError{http.StatusUnprocessableEntity, "UNDEFINED_RESULT"}},
	{calculator.ErrNonFiniteResult, apiError{http.StatusUnprocessableEntity, "NON_FINITE_RESULT"}},
}

// writeError sends err as a JSON error response with the right status code.
// Unknown errors become a generic 500 so internal details never leak.
func writeError(w http.ResponseWriter, err error) {
	for _, m := range errorMappings {
		if errors.Is(err, m.target) {
			writeJSON(w, m.api.status, ErrorResponse{
				Error: ErrorBody{Code: m.api.code, Message: err.Error()},
			})
			return
		}
	}

	slog.Error("unexpected error", "err", err)
	writeJSON(w, http.StatusInternalServerError, ErrorResponse{
		Error: ErrorBody{Code: "INTERNAL_ERROR", Message: "internal server error"},
	})
}

// writeJSON sets the JSON content type, writes the status code and encodes v.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("encode response", "err", err)
	}
}