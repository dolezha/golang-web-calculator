package test

import (
	"bytes"
	"calculator/handlers"
	"calculator/models"
	"calculator/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateHandler(t *testing.T) {
	tests := []struct {
		name           string
		body           models.RequestBody
		expectedStatus int
		expectedResult *float64
		expectedError  *string
	}{
		{
			name:           "Valid expression",
			body:           models.RequestBody{Expression: "2 + 2"},
			expectedStatus: http.StatusOK,
			expectedResult: utils.ToPtr(4.0),
			expectedError:  nil,
		},
		{
			name:           "Division by zero",
			body:           models.RequestBody{Expression: "2 / 0"},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedResult: nil,
			expectedError:  utils.ToPtr("Division by zero"),
		},
		{
			name:           "Invalid expression",
			body:           models.RequestBody{Expression: "2 + + 2"},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedResult: nil,
			expectedError:  utils.ToPtr("Expression is not valid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(handlers.CalculateHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, rr.Code)
			}

			var response models.ResponseBody
			if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
				t.Fatal(err)
			}

			if tt.expectedResult != nil && *response.Result != *tt.expectedResult {
				t.Errorf("expected result %v, got %v", *tt.expectedResult, *response.Result)
			}

			if tt.expectedError != nil && *response.Error != *tt.expectedError {
				t.Errorf("expected error %v, got %v", *tt.expectedError, *response.Error)
			}
		})
	}
}
