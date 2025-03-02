package handlers

import (
	"bytes"
	"calculator/services"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTaskHandler(t *testing.T) {
	// Создаем тестовое выражение для генерации задач
	_, _ = services.CreateExpression("2+2")

	t.Run("get task", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/internal/task", nil)
		w := httptest.NewRecorder()

		TaskHandler(w, req)

		if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
			t.Errorf("TaskHandler() status = %v, want %v or %v", w.Code, http.StatusOK, http.StatusNotFound)
		}

		if w.Code == http.StatusOK {
			var response struct {
				Task struct {
					ID string `json:"id"`
				} `json:"task"`
			}
			json.NewDecoder(w.Body).Decode(&response)

			if response.Task.ID == "" {
				t.Error("TaskHandler() response should contain task ID")
			}
		}
	})

	t.Run("submit result", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"id":     "test_task",
			"result": 4.0,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/internal/task", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		TaskHandler(w, req)

		if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
			t.Errorf("TaskHandler() status = %v, want %v or %v", w.Code, http.StatusOK, http.StatusNotFound)
		}
	})
}
