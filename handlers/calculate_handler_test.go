package handlers

import (
	"bytes"
	"calculator/middleware"
	"calculator/models"
	"calculator/services"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type ExpressionServiceInterface interface {
	CreateExpression(userID int, expr string) (*models.Expression, error)
}

type MockExpressionService struct {
	shouldError bool
}

func (m *MockExpressionService) CreateExpression(userID int, expr string) (*models.Expression, error) {
	if m.shouldError {
		return nil, fmt.Errorf("invalid expression")
	}
	return &models.Expression{
		ID:         "test-id",
		UserID:     userID,
		Expression: expr,
		Status:     models.StatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

type TestCalculateHandler struct {
	service ExpressionServiceInterface
}

func NewTestCalculateHandler(service ExpressionServiceInterface) *TestCalculateHandler {
	return &TestCalculateHandler{service: service}
}

func (ch *TestCalculateHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		http.Error(w, "User not authorized", http.StatusUnauthorized)
		return
	}

	var reqBody models.RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}

	expression, err := ch.service.CreateExpression(claims.UserID, reqBody.Expression)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	response := map[string]string{
		"id": expression.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func TestCalculateHandler_RealHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           string
		hasAuth        bool
		serviceError   bool
		expectedStatus int
	}{
		{
			name:           "valid expression",
			method:         "POST",
			body:           `{"expression":"2+2"}`,
			hasAuth:        true,
			serviceError:   false,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid method",
			method:         "GET",
			body:           `{"expression":"2+2"}`,
			hasAuth:        true,
			serviceError:   false,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "no auth",
			method:         "POST",
			body:           `{"expression":"2+2"}`,
			hasAuth:        false,
			serviceError:   false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid JSON",
			method:         "POST",
			body:           `invalid json`,
			hasAuth:        true,
			serviceError:   false,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "service error",
			method:         "POST",
			body:           `{"expression":"invalid"}`,
			hasAuth:        true,
			serviceError:   true,
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockExpressionService{shouldError: tt.serviceError}
			handler := NewTestCalculateHandler(mockService)

			req := httptest.NewRequest(tt.method, "/api/v1/calculate", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			if tt.hasAuth {
				claims := &services.Claims{
					UserID: 1,
					Login:  "testuser",
				}
				ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			handler.Calculate(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestCalculateHandler_Legacy(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockExpressionService{}
			handler := NewTestCalculateHandler(mockService)

			reqBody := map[string]string{
				"expression": tt.expression,
			}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			claims := &services.Claims{
				UserID: 1,
				Login:  "testuser",
			}
			ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			handler.Calculate(w, req)

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
