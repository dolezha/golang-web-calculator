package handlers

import (
	"calculator/models"
	"calculator/services"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExpressionHandler(t *testing.T) {
	// Создаем тестовое выражение
	expr, _ := services.CreateExpression("2+2")

	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{
			name:       "existing expression",
			path:       "/api/v1/expressions/" + expr.ID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-existing expression",
			path:       "/api/v1/expressions/nonexistent",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			ExpressionHandler(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("ExpressionHandler() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var response struct {
					Expression struct {
						ID     string                  `json:"id"`
						Status models.ExpressionStatus `json:"status"`
					} `json:"expression"`
				}
				json.NewDecoder(w.Body).Decode(&response)

				if response.Expression.ID != expr.ID {
					t.Errorf("ExpressionHandler() id = %v, want %v", response.Expression.ID, expr.ID)
				}
			}
		})
	}
}
