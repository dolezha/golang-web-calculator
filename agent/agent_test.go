package agent

import (
	"calculator/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGetEnvInt(t *testing.T) {
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")

	result := getEnvInt("TEST_INT", 10)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}

	result = getEnvInt("NON_EXISTING", 10)
	if result != 10 {
		t.Errorf("Expected 10, got %d", result)
	}

	os.Setenv("INVALID_INT", "not-a-number")
	defer os.Unsetenv("INVALID_INT")

	result = getEnvInt("INVALID_INT", 10)
	if result != 10 {
		t.Errorf("Expected 10, got %d", result)
	}
}

func TestGetServerURL(t *testing.T) {
	os.Setenv("SERVER_URL", "http://test:8080")
	defer os.Unsetenv("SERVER_URL")

	result := getServerURL()
	if result != "http://test:8080" {
		t.Errorf("Expected http://test:8080, got %s", result)
	}

	os.Unsetenv("SERVER_URL")
	result = getServerURL()
	if result != "http://calc-service:8080" {
		t.Errorf("Expected http://calc-service:8080, got %s", result)
	}
}

func TestCompute(t *testing.T) {
	tests := []struct {
		name     string
		task     *models.Task
		expected float64
	}{
		{
			name: "addition",
			task: &models.Task{
				Arg1:      "2",
				Arg2:      "3",
				Operation: "+",
			},
			expected: 5,
		},
		{
			name: "subtraction",
			task: &models.Task{
				Arg1:      "5",
				Arg2:      "3",
				Operation: "-",
			},
			expected: 2,
		},
		{
			name: "multiplication",
			task: &models.Task{
				Arg1:      "4",
				Arg2:      "3",
				Operation: "*",
			},
			expected: 12,
		},
		{
			name: "division",
			task: &models.Task{
				Arg1:      "6",
				Arg2:      "2",
				Operation: "/",
			},
			expected: 3,
		},
		{
			name: "division by zero",
			task: &models.Task{
				Arg1:      "6",
				Arg2:      "0",
				Operation: "/",
			},
			expected: 0,
		},
		{
			name: "unknown operation",
			task: &models.Task{
				Arg1:      "6",
				Arg2:      "2",
				Operation: "^",
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compute(tt.task)
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestGetTaskFromServer(t *testing.T) {
	task := &models.Task{
		ID:        "test-task",
		Operation: "+",
		Arg1:      "2",
		Arg2:      "3",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/internal/task" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}))
	defer server.Close()

	originalURL := serverURL
	serverURL = server.URL
	defer func() { serverURL = originalURL }()

	result, err := getTaskFromServer()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.ID != task.ID {
		t.Errorf("Expected task ID %s, got %s", task.ID, result.ID)
	}
}

func TestGetTaskFromServer_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	originalURL := serverURL
	serverURL = server.URL
	defer func() { serverURL = originalURL }()

	_, err := getTaskFromServer()
	if err == nil {
		t.Error("Expected error for 404 response")
	}
}

func TestSubmitTaskResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/internal/task/test-task" && r.Method == "POST" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer server.Close()

	originalURL := serverURL
	serverURL = server.URL
	defer func() { serverURL = originalURL }()

	err := submitTaskResult("test-task", 42.0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetTaskResult(t *testing.T) {
	result := 42.0
	task := &models.Task{
		ID:     "test-task",
		Result: &result,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/internal/task/test-task" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	originalURL := serverURL
	serverURL = server.URL
	defer func() { serverURL = originalURL }()

	resultTask, err := getTaskResult("test-task")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resultTask.ID != task.ID {
		t.Errorf("Expected task ID %s, got %s", task.ID, resultTask.ID)
	}

	if resultTask.Result == nil || *resultTask.Result != *task.Result {
		t.Errorf("Expected result %v, got %v", task.Result, resultTask.Result)
	}
}

func TestComputeWithDependency(t *testing.T) {
	dependencyResult := 5.0
	dependencyTask := &models.Task{
		ID:     "dep-task",
		Result: &dependencyResult,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/internal/task/dep-task" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(dependencyTask)
		}
	}))
	defer server.Close()

	originalURL := serverURL
	serverURL = server.URL
	defer func() { serverURL = originalURL }()

	task := &models.Task{
		Arg1:      "$dep-task",
		Arg2:      "3",
		Operation: "+",
	}

	result := compute(task)
	if result != 8.0 {
		t.Errorf("Expected 8.0, got %f", result)
	}
}

func TestStartAgent(t *testing.T) {

	originalComputingPower := computingPower
	originalServerURL := serverURL

	computingPower = 1
	serverURL = "http://test:8080"

	defer func() {
		computingPower = originalComputingPower
		serverURL = originalServerURL
	}()

	if os.Getenv("RUN_START_AGENT_TEST") == "1" {
		go StartAgent()
		time.Sleep(10 * time.Millisecond)
	}
}

func TestWorker(t *testing.T) {
	taskCalled := false
	resultCalled := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/internal/task" && r.Method == "GET" {
			taskCalled = true
			if !resultCalled {
				task := &models.Task{
					ID:            "test-task",
					Arg1:          "2",
					Arg2:          "3",
					Operation:     "+",
					OperationTime: 1,
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(task)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else if strings.HasPrefix(r.URL.Path, "/internal/task/") && r.Method == "POST" {
			resultCalled = true
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	originalURL := serverURL
	serverURL = server.URL
	defer func() { serverURL = originalURL }()

	done := make(chan bool)
	go func() {
		worker(0)
		done <- true
	}()

	time.Sleep(50 * time.Millisecond)

	if !taskCalled {
		t.Error("Worker should have called getTaskFromServer")
	}
	if !resultCalled {
		t.Error("Worker should have called submitTaskResult")
	}
}

func TestComputeWithInvalidArg(t *testing.T) {
	task := &models.Task{
		Arg1:      "invalid",
		Arg2:      "3",
		Operation: "+",
	}

	result := compute(task)
	if result != 3.0 {
		t.Errorf("Expected 3.0 (0 + 3), got %f", result)
	}
}

func TestSubmitTaskResult_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	}))
	defer server.Close()

	originalURL := serverURL
	serverURL = server.URL
	defer func() { serverURL = originalURL }()

	err := submitTaskResult("test-task", 42.0)
	if err == nil {
		t.Error("Expected error for server error response")
	}
}

func TestGetTaskFromServer_UnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	originalURL := serverURL
	serverURL = server.URL
	defer func() { serverURL = originalURL }()

	_, err := getTaskFromServer()
	if err == nil {
		t.Error("Expected error for unexpected status code")
	}
}

func TestGetTaskFromServer_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	originalURL := serverURL
	serverURL = server.URL
	defer func() { serverURL = originalURL }()

	_, err := getTaskFromServer()
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestGetTaskResult_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	originalURL := serverURL
	serverURL = server.URL
	defer func() { serverURL = originalURL }()

	_, err := getTaskResult("nonexistent-task")
	if err == nil {
		t.Error("Expected error for task not found")
	}
}

func TestGetTaskResult_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	originalURL := serverURL
	serverURL = server.URL
	defer func() { serverURL = originalURL }()

	_, err := getTaskResult("test-task")
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestSubmitTaskResult_HTTPError(t *testing.T) {
	originalURL := serverURL
	serverURL = "http://invalid-url-that-does-not-exist:99999"
	defer func() { serverURL = originalURL }()

	err := submitTaskResult("test-task", 42.0)
	if err == nil {
		t.Error("Expected error for HTTP request failure")
	}
}

func TestGetTaskFromServer_HTTPError(t *testing.T) {
	originalURL := serverURL
	serverURL = "http://invalid-url-that-does-not-exist:99999"
	defer func() { serverURL = originalURL }()

	_, err := getTaskFromServer()
	if err == nil {
		t.Error("Expected error for HTTP request failure")
	}
}

func TestGetTaskResult_HTTPError(t *testing.T) {
	originalURL := serverURL
	serverURL = "http://invalid-url-that-does-not-exist:99999"
	defer func() { serverURL = originalURL }()

	_, err := getTaskResult("test-task")
	if err == nil {
		t.Error("Expected error for HTTP request failure")
	}
}

func TestWorker_ErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	originalURL := serverURL
	serverURL = server.URL
	defer func() { serverURL = originalURL }()

	done := make(chan bool)
	go func() {
		worker(0)
		done <- true
	}()

	time.Sleep(50 * time.Millisecond)

}

func TestComputeWithDependencyWaiting(t *testing.T) {
	callCount := 0
	dependencyResult := 5.0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/internal/task/dep-task" {
			callCount++
			w.Header().Set("Content-Type", "application/json")

			if callCount == 1 {
				task := &models.Task{
					ID:     "dep-task",
					Result: nil,
				}
				json.NewEncoder(w).Encode(task)
			} else {
				task := &models.Task{
					ID:     "dep-task",
					Result: &dependencyResult,
				}
				json.NewEncoder(w).Encode(task)
			}
		}
	}))
	defer server.Close()

	originalURL := serverURL
	serverURL = server.URL
	defer func() { serverURL = originalURL }()

	task := &models.Task{
		Arg1:      "$dep-task",
		Arg2:      "3",
		Operation: "+",
	}

	result := compute(task)
	if result != 8.0 {
		t.Errorf("Expected 8.0, got %f", result)
	}

	if callCount < 2 {
		t.Errorf("Expected at least 2 calls to dependency task, got %d", callCount)
	}
}

func TestComputeWithDependencyError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/internal/task/dep-task" {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	originalURL := serverURL
	serverURL = server.URL
	defer func() { serverURL = originalURL }()

	task := &models.Task{
		Arg1:      "$dep-task",
		Arg2:      "3",
		Operation: "+",
	}

	done := make(chan float64)
	go func() {
		result := compute(task)
		done <- result
	}()

	select {
	case result := <-done:
		t.Errorf("Expected compute to keep waiting, but got result: %f", result)
	case <-time.After(200 * time.Millisecond):
	}
}

func TestWorker_SubmitError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/internal/task" && r.Method == "GET" {
			task := &models.Task{
				ID:            "test-task",
				Arg1:          "2",
				Arg2:          "3",
				Operation:     "+",
				OperationTime: 1,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		} else if strings.HasPrefix(r.URL.Path, "/internal/task/") && r.Method == "POST" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Server error"))
		}
	}))
	defer server.Close()

	originalURL := serverURL
	serverURL = server.URL
	defer func() { serverURL = originalURL }()

	done := make(chan bool)
	go func() {
		worker(0)
		done <- true
	}()

	time.Sleep(100 * time.Millisecond)

}

func TestGlobalVariables(t *testing.T) {
	if computingPower <= 0 {
		t.Error("computingPower should be positive")
	}

	if serverURL == "" {
		t.Error("serverURL should not be empty")
	}

	originalPower := computingPower
	originalURL := serverURL

	computingPower = 10
	serverURL = "http://test:9999"

	if computingPower != 10 {
		t.Error("Should be able to modify computingPower")
	}

	if serverURL != "http://test:9999" {
		t.Error("Should be able to modify serverURL")
	}

	computingPower = originalPower
	serverURL = originalURL
}

func TestStartAgentPrintStatement(t *testing.T) {
	originalPower := computingPower
	computingPower = 2
	defer func() { computingPower = originalPower }()

	if os.Getenv("TEST_START_AGENT_PRINT") == "1" {
		go StartAgent()
		time.Sleep(10 * time.Millisecond)
	}
}

func TestStartAgentSelectBlock(t *testing.T) {
	originalPower := computingPower
	computingPower = 1
	defer func() { computingPower = originalPower }()

	done := make(chan bool)
	go func() {
		defer func() {
			if r := recover(); r != nil {
			}
			done <- true
		}()

		go StartAgent()

		time.Sleep(50 * time.Millisecond)

		panic("test exit")
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Error("StartAgent test timed out")
	}
}

func TestComputeAllOperations(t *testing.T) {
	tests := []struct {
		arg1      string
		arg2      string
		operation string
		expected  float64
	}{
		{"5", "3", "+", 8},
		{"5", "3", "-", 2},
		{"5", "3", "*", 15},
		{"6", "3", "/", 2},
		{"5", "0", "/", 0},
		{"5", "3", "%", 0},
	}

	for _, tt := range tests {
		task := &models.Task{
			Arg1:      tt.arg1,
			Arg2:      tt.arg2,
			Operation: tt.operation,
		}
		result := compute(task)
		if result != tt.expected {
			t.Errorf("compute(%s %s %s) = %f, expected %f", tt.arg1, tt.operation, tt.arg2, result, tt.expected)
		}
	}
}
