package test

import (
	"calculator/services"
	"testing"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		name           string
		expression     string
		expectedResult float64
		expectedError  bool
	}{
		{
			name:           "Valid expression 1",
			expression:     "2 + 2",
			expectedResult: 4.0,
			expectedError:  false,
		},
		{
			name:           "Division by zero",
			expression:     "2 / 0",
			expectedResult: 0,
			expectedError:  true,
		},
		{
			name:           "Invalid expression",
			expression:     "2 + + 2",
			expectedResult: 0,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := services.Calc(tt.expression)

			if tt.expectedError && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("did not expect error, got %v", err)
			}
			if !tt.expectedError && result != tt.expectedResult {
				t.Errorf("expected result %v, got %v", tt.expectedResult, result)
			}
		})
	}
}
