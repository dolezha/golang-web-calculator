package test

import (
	"calculator/utils"
	"testing"
)

func TestToPtr(t *testing.T) {
	tests := []struct {
		name  string
		value int
	}{
		{
			name:  "Integer to pointer",
			value: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := utils.ToPtr(tt.value)

			if *ptr != tt.value {
				t.Errorf("expected pointer value %v, got %v", tt.value, *ptr)
			}
		})
	}
}
