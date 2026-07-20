package main

import "testing"

type testCase struct {
	name     string
	a, b     int
	expected int
	err      bool
}

func TestAll(t *testing.T) {
	tests := []testCase{
		// Add
		{"Add positive", 2, 3, 5, false},
		{"Add zero", 0, 5, 5, false},
		{"Add negative", -2, -3, -5, false},
		{"Add mixed", -2, 3, 1, false},

		// Multiply
		{"Multiply positive", 2, 3, 6, false},
		{"Multiply zero", 5, 0, 0, false},
		{"Multiply negative", -2, 3, -6, false},
		{"Multiply mixed", -2, -3, 6, false},

		// Divide
		{"Divide positive", 6, 3, 2, false},
		{"Divide negative", -6, 3, -2, false},
		{"Divide by one", 7, 1, 7, false},
		{"Divide by zero", 10, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result int
			var err error

			// Определяем, какую функцию вызывать
			switch {
			case tt.name[:3] == "Add":
				result = Add(tt.a, tt.b)
			case tt.name[:8] == "Multiply":
				result = Multiply(tt.a, tt.b)
			case tt.name[:6] == "Divide":
				result, err = Divide(tt.a, tt.b)
			}

			if tt.err && err == nil {
				t.Errorf("%s expected error, got nil", tt.name)
			}
			if !tt.err && err != nil {
				t.Errorf("%s error: %v", tt.name, err)
			}
			if result != tt.expected {
				t.Errorf("%s = %d; want %d", tt.name, result, tt.expected)
			}
		})
	}
}
