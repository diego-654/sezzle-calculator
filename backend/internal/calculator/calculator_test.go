package calculator

import "testing"

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{name: "dos positivos", a: 2, b: 3, want: 5},
		{name: "con negativo", a: -4, b: 1, want: -3},
		{name: "con cero", a: 7, b: 0, want: 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Add(tt.a, tt.b)
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if got != tt.want {
				t.Errorf("Add(%v, %v) = %v, se esperaba %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{name: "resultado positivo", a: 10, b: 4, want: 6},
		{name: "resultado negativo", a: 4, b: 10, want: -6},
		{name: "restar un negativo", a: 5, b: -3, want: 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Subtract(tt.a, tt.b)
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if got != tt.want {
				t.Errorf("Subtract(%v, %v) = %v, se esperaba %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{name: "dos positivos", a: 3, b: 4, want: 12},
		{name: "positivo por negativo", a: 3, b: -4, want: -12},
		{name: "por cero", a: 99, b: 0, want: 0},
		{name: "decimales", a: 2.5, b: 4, want: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Multiply(tt.a, tt.b)
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if got != tt.want {
				t.Errorf("Multiply(%v, %v) = %v, se esperaba %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{name: "division exacta", a: 10, b: 2, want: 5},
		{name: "resultado decimal", a: 1, b: 4, want: 0.25},
		{name: "divisor negativo", a: 9, b: -3, want: -3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if got != tt.want {
				t.Errorf("Divide(%v, %v) = %v, se esperaba %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestPower(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{name: "exponente positivo", a: 2, b: 10, want: 1024},
		{name: "exponente cero", a: 5, b: 0, want: 1},
		{name: "exponente negativo", a: 2, b: -1, want: 0.5},
		{name: "base negativa con exponente entero", a: -2, b: 3, want: -8},
		{name: "exponente fraccionario", a: 9, b: 0.5, want: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Power(tt.a, tt.b)
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if got != tt.want {
				t.Errorf("Power(%v, %v) = %v, se esperaba %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		want float64
	}{
		{name: "cuadrado perfecto", a: 16, want: 4},
		{name: "cero", a: 0, want: 0},
		{name: "decimal", a: 2.25, want: 1.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sqrt(tt.a)
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if got != tt.want {
				t.Errorf("Sqrt(%v) = %v, se esperaba %v", tt.a, got, tt.want)
			}
		})
	}
}

func TestPercent(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{name: "15% de 200", a: 200, b: 15, want: 30},
		{name: "100% de 50", a: 50, b: 100, want: 50},
		{name: "0% de 80", a: 80, b: 0, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Percent(tt.a, tt.b)
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if got != tt.want {
				t.Errorf("Percent(%v, %v) = %v, se esperaba %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
