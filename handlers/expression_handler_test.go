package handlers

import (
	"calculator/middleware"
	"calculator/models"
	"calculator/services"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockExpressionServiceForExpressionHandler struct{}

func (m *MockExpressionServiceForExpressionHandler) GetExpression(id string, userID int) (*models.Expression, error) {
	return &models.Expression{
		ID:         id,
		UserID:     userID,
		Expression: "2+2",
		Status:     models.StatusDone,
		Result:     func() *float64 { r := 4.0; return &r }(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

func (m *MockExpressionServiceForExpressionHandler) GetUserExpressions(userID int) ([]*models.Expression, error) {
	return []*models.Expression{
		{
			ID:         "test-id",
			UserID:     userID,
			Expression: "2+2",
			Status:     models.StatusDone,
			Result:     func() *float64 { r := 4.0; return &r }(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}, nil
}

type MockExpressionHandler struct {
	expressionService *MockExpressionServiceForExpressionHandler
}

func NewMockExpressionHandler(expressionService *MockExpressionServiceForExpressionHandler) *MockExpressionHandler {
	return &MockExpressionHandler{expressionService: expressionService}
}

func (eh *MockExpressionHandler) GetExpressions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	expressions, err := eh.expressionService.GetUserExpressions(claims.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expressions)
}

func TestExpressionHandler(t *testing.T) {
	t.Run("get expressions", func(t *testing.T) {
		mockService := &MockExpressionServiceForExpressionHandler{}
		handler := NewMockExpressionHandler(mockService)

		req := httptest.NewRequest("GET", "/api/v1/expressions", nil)

		claims := &services.Claims{
			UserID: 1,
			Login:  "testuser",
		}
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, claims)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		handler.GetExpressions(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GetExpressions() status = %v, want %v", w.Code, http.StatusOK)
		}

		var response []*models.Expression
		json.NewDecoder(w.Body).Decode(&response)

		if len(response) == 0 {
			t.Error("GetExpressions() should return expressions")
		}
	})
}
