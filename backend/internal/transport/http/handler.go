package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/diego-654/sezzle-calculator/backend/internal/calculator"
)

// maxBodyBytes caps the request body. Two numbers fit easily in 1 KB,
// so anything bigger is rejected before it is read into memory.
const maxBodyBytes = 1 << 10

type operation struct {
	needsB bool                                // false for unary operations like sqrt
	fn     func(a, b float64) (float64, error) // the domain function to call
}

// operations maps the {op} path segment to its domain function.
var operations = map[string]operation{
	"add":      {needsB: true, fn: calculator.Add},
	"subtract": {needsB: true, fn: calculator.Subtract},
	"multiply": {needsB: true, fn: calculator.Multiply},
	"divide":   {needsB: true, fn: calculator.Divide},
	"power":    {needsB: true, fn: calculator.Power},
	"percent":  {needsB: true, fn: calculator.Percent},
	"sqrt": {needsB: false, fn: func(a, _ float64) (float64, error) {
		return calculator.Sqrt(a)
	}},
}

// handleCalculate serves POST /api/v1/{op} for every operation.
func handleCalculate(w http.ResponseWriter, r *http.Request) {
	op, ok := operations[r.PathValue("op")]
	if !ok {
		writeError(w, calculator.ErrInvalidOperation)
		return
	}

	req, err := decodeRequest(w, r)
	if err != nil {
		writeError(w, err)
		return
	}

	if req.A == nil {
		writeError(w, fmt.Errorf("%w: missing operand \"a\"", calculator.ErrInvalidInput))
		return
	}
	var b float64
	if op.needsB {
		if req.B == nil {
			writeError(w, fmt.Errorf("%w: missing operand \"b\"", calculator.ErrInvalidInput))
			return
		}
		b = *req.B
	}

	result, err := op.fn(*req.A, b)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, CalculateResponse{Result: result})
}

// decodeRequest reads exactly one JSON object from the body and rejects
// oversized bodies, unknown fields and trailing data.
func decodeRequest(w http.ResponseWriter, r *http.Request) (CalculateRequest, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req CalculateRequest
	if err := dec.Decode(&req); err != nil {
		return req, fmt.Errorf("%w: invalid JSON body", calculator.ErrInvalidInput)
	}
	// A second Decode must hit EOF; otherwise the body had extra data.
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return req, fmt.Errorf("%w: body must contain a single JSON object", calculator.ErrInvalidInput)
	}
	return req, nil
}
