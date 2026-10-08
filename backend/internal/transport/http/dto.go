// Package httpapi exposes the calculator over HTTP. It only translates
// JSON into domain calls and domain errors into HTTP responses.
package httpapi

// CalculateRequest is the body of POST /api/v1/{op}.
type CalculateRequest struct {
	A *float64 `json:"a"`
	B *float64 `json:"b"` 
}

// CalculateResponse is a successful response: {"result": 5}.
type CalculateResponse struct {
	Result float64 `json:"result"`
}

// ErrorBody describes an error: a stable code for machines and a
// readable message for people.
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse wraps the error: {"error": {"code": "...", "message": "..."}}.
type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}