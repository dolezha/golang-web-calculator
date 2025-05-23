package handlers

import (
	"calculator/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockExpressionServiceForTaskHandler struct{}

func (m *MockExpressionServiceForTaskHandler) GetNextTask() (*models.Task, error) {
	return &models.Task{
		ID:            "test-task-id",
		ExpressionID:  "test-expr-id",
		Arg1:          "2",
		Arg2:          "3",
		Operation:     "+",
		OperationTime: 1000,
		Status:        "pending",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

func (m *MockExpressionServiceForTaskHandler) GetTaskByID(taskID string) (*models.Task, error) {
	return &models.Task{
		ID:            taskID,
		ExpressionID:  "test-expr-id",
		Arg1:          "2",
		Arg2:          "3",
		Operation:     "+",
		OperationTime: 1000,
		Status:        "computing",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

type MockTaskHandler struct {
	expressionService *MockExpressionServiceForTaskHandler
}

func NewMockTaskHandler(expressionService *MockExpressionServiceForTaskHandler) *MockTaskHandler {
	return &MockTaskHandler{expressionService: expressionService}
}

func (th *MockTaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	task, err := th.expressionService.GetNextTask()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func TestTaskHandler(t *testing.T) {
	t.Run("get task", func(t *testing.T) {
		mockService := &MockExpressionServiceForTaskHandler{}
		handler := NewMockTaskHandler(mockService)

		req := httptest.NewRequest("GET", "/internal/task", nil)
		w := httptest.NewRecorder()

		handler.GetTask(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GetTask() status = %v, want %v", w.Code, http.StatusOK)
		}

		var response models.Task
		json.NewDecoder(w.Body).Decode(&response)

		if response.ID == "" {
			t.Error("GetTask() should return task with ID")
		}
	})
}
