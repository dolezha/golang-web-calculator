package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateHandler(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		wantStatus int
		wantID     bool
	}{
		{
			name:       "valid expression",
			expression: "2+2*2",
			wantStatus: http.StatusCreated,
			wantID:     true,
		},
		{
			name:       "invalid expression",
			expression: "2++2",
			wantStatus: http.StatusUnprocessableEntity,
			wantID:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := map[string]string{
				"expression": tt.expression,
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			CalculateHandler(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("CalculateHandler() status = %v, want %v", w.Code, tt.wantStatus)
			}

			var response map[string]string
			json.NewDecoder(w.Body).Decode(&response)

			if tt.wantID && response["id"] == "" {
				t.Error("CalculateHandler() response should contain ID")
			}
		})
	}
}
