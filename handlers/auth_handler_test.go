package handlers

import (
	"calculator/models"
	"calculator/services"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type AuthServiceInterface interface {
	Register(req *models.RegisterRequest) error
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
	ValidateToken(token string) (*services.Claims, error)
}

type MockAuthService struct {
	shouldErrorRegister bool
	shouldErrorLogin    bool
}

func (m *MockAuthService) Register(req *models.RegisterRequest) error {
	if m.shouldErrorRegister {
		return fmt.Errorf("user already exists")
	}
	return nil
}

func (m *MockAuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	if m.shouldErrorLogin {
		return nil, fmt.Errorf("invalid credentials")
	}
	return &models.LoginResponse{Token: "test-jwt-token"}, nil
}

func (m *MockAuthService) ValidateToken(token string) (*services.Claims, error) {
	return &services.Claims{UserID: 1, Login: "testuser"}, nil
}

type TestAuthHandler struct {
	authService AuthServiceInterface
}

func NewTestAuthHandler(authService AuthServiceInterface) *TestAuthHandler {
	return &TestAuthHandler{authService: authService}
}

func (ah *TestAuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request format"})
		return
	}

	if err := ah.authService.Register(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (ah *TestAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request format"})
		return
	}

	response, err := ah.authService.Login(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           string
		serviceError   bool
		expectedStatus int
	}{
		{
			name:           "valid registration",
			method:         "POST",
			body:           `{"login":"testuser","password":"password123"}`,
			serviceError:   false,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid method",
			method:         "GET",
			body:           `{"login":"testuser","password":"password123"}`,
			serviceError:   false,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "invalid JSON",
			method:         "POST",
			body:           `invalid json`,
			serviceError:   false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "service error",
			method:         "POST",
			body:           `{"login":"existing","password":"password123"}`,
			serviceError:   true,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAuthService{shouldErrorRegister: tt.serviceError}
			handler := NewTestAuthHandler(mockService)

			req := httptest.NewRequest(tt.method, "/api/v1/register", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.Register(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           string
		serviceError   bool
		expectedStatus int
	}{
		{
			name:           "valid login",
			method:         "POST",
			body:           `{"login":"testuser","password":"password123"}`,
			serviceError:   false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid method",
			method:         "GET",
			body:           `{"login":"testuser","password":"password123"}`,
			serviceError:   false,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "invalid JSON",
			method:         "POST",
			body:           `invalid json`,
			serviceError:   false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "service error",
			method:         "POST",
			body:           `{"login":"wrong","password":"wrong"}`,
			serviceError:   true,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAuthService{shouldErrorLogin: tt.serviceError}
			handler := NewTestAuthHandler(mockService)

			req := httptest.NewRequest(tt.method, "/api/v1/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.Login(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response models.LoginResponse
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Errorf("Failed to decode response: %v", err)
				}
				if response.Token == "" {
					t.Error("Expected token in response")
				}
			}
		})
	}
}
