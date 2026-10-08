package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantResult float64 // checked only when wantStatus is 200
		wantCode   string  // checked only when set
	}{
		// Successful operations
		{name: "add", method: "POST", path: "/api/v1/add", body: `{"a":2,"b":3}`, wantStatus: 200, wantResult: 5},
		{name: "subtract", method: "POST", path: "/api/v1/subtract", body: `{"a":10,"b":4}`, wantStatus: 200, wantResult: 6},
		{name: "multiply", method: "POST", path: "/api/v1/multiply", body: `{"a":3,"b":4}`, wantStatus: 200, wantResult: 12},
		{name: "divide", method: "POST", path: "/api/v1/divide", body: `{"a":10,"b":4}`, wantStatus: 200, wantResult: 2.5},
		{name: "power", method: "POST", path: "/api/v1/power", body: `{"a":2,"b":10}`, wantStatus: 200, wantResult: 1024},
		{name: "percent", method: "POST", path: "/api/v1/percent", body: `{"a":200,"b":15}`, wantStatus: 200, wantResult: 30},
		{name: "sqrt without b", method: "POST", path: "/api/v1/sqrt", body: `{"a":9}`, wantStatus: 200, wantResult: 3},

		// Domain errors: 422
		{name: "divide by zero", method: "POST", path: "/api/v1/divide", body: `{"a":1,"b":0}`, wantStatus: 422, wantCode: "DIVISION_BY_ZERO"},
		{name: "sqrt of negative", method: "POST", path: "/api/v1/sqrt", body: `{"a":-4}`, wantStatus: 422, wantCode: "NEGATIVE_SQUARE_ROOT"},

		// Bad input: 400
		{name: "missing b", method: "POST", path: "/api/v1/add", body: `{"a":1}`, wantStatus: 400, wantCode: "INVALID_INPUT"},
		{name: "missing a", method: "POST", path: "/api/v1/sqrt", body: `{}`, wantStatus: 400, wantCode: "INVALID_INPUT"},
		{name: "invalid JSON", method: "POST", path: "/api/v1/add", body: `nope`, wantStatus: 400, wantCode: "INVALID_INPUT"},
		{name: "unknown field", method: "POST", path: "/api/v1/add", body: `{"a":1,"b":2,"c":3}`, wantStatus: 400, wantCode: "INVALID_INPUT"},
		{name: "two JSON objects", method: "POST", path: "/api/v1/add", body: `{"a":1,"b":2}{}`, wantStatus: 400, wantCode: "INVALID_INPUT"},

		// Routing
		{name: "unknown operation", method: "POST", path: "/api/v1/modulo", body: `{"a":1,"b":2}`, wantStatus: 404, wantCode: "INVALID_OPERATION"},
		{name: "wrong method", method: "GET", path: "/api/v1/add", wantStatus: 405},
	}

	router := NewRouter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body.String())
			}

			if tt.wantStatus == http.StatusOK {
				var resp CalculateResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp.Result != tt.wantResult {
					t.Errorf("result = %v, want %v", resp.Result, tt.wantResult)
				}
			}

			if tt.wantCode != "" {
				var resp ErrorResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("decode error response: %v", err)
				}
				if resp.Error.Code != tt.wantCode {
					t.Errorf("error code = %q, want %q", resp.Error.Code, tt.wantCode)
				}
			}
		})
	}
}

func TestHealth(t *testing.T) {
	rec := httptest.NewRecorder()
	NewRouter().ServeHTTP(rec, httptest.NewRequest("GET", "/health", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}
