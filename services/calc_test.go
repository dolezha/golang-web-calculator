package services

import "testing"

func TestCalc(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		want    float64
		wantErr bool
	}{
		{
			name:    "simple addition",
			expr:    "2+2",
			want:    4,
			wantErr: false,
		},
		{
			name:    "complex expression",
			expr:    "2+2*2",
			want:    6,
			wantErr: false,
		},
		{
			name:    "invalid expression",
			expr:    "2++2",
			want:    0,
			wantErr: true,
		},
		{
			name:    "division by zero",
			expr:    "1/0",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calc(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Calc() = %v, want %v", got, tt.want)
			}
		})
	}
}
