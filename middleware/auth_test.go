package middleware

import (
	"calculator/models"
	"calculator/services"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MockDatabase struct{}

func (m *MockDatabase) CreateUser(login, passwordHash string) (*models.User, error) {
	return nil, nil
}

func (m *MockDatabase) GetUserByLogin(login string) (*models.User, error) {
	return nil, fmt.Errorf("user not found")
}

func (m *MockDatabase) Close() error {
	return nil
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	mockDB := &MockDatabase{}
	authService := services.NewAuthService(mockDB, "test-secret")
	middleware := AuthMiddleware(authService)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	protectedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	protectedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	mockDB := &MockDatabase{}
	authService := services.NewAuthService(mockDB, "test-secret")
	middleware := AuthMiddleware(authService)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	protectedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "invalid-format")
	rr := httptest.NewRecorder()
	protectedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	mockDB := &MockDatabase{}
	authService := services.NewAuthService(mockDB, "test-secret")
	middleware := AuthMiddleware(authService)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	protectedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rr := httptest.NewRecorder()
	protectedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	mockDB := &MockDatabase{}
	authService := services.NewAuthService(mockDB, "test-secret")
	middleware := AuthMiddleware(authService)

	claims := &services.Claims{
		UserID: 1,
		Login:  "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		user, ok := GetUserFromContext(r)
		if !ok {
			t.Error("Expected user in context")
		}
		if user.UserID != 1 {
			t.Errorf("Expected UserID 1, got %d", user.UserID)
		}
		w.WriteHeader(http.StatusOK)
	})

	protectedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rr := httptest.NewRecorder()
	protectedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if !handlerCalled {
		t.Error("Expected handler to be called")
	}
}

func TestGetUserFromContext(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)

	_, ok := GetUserFromContext(req)
	if ok {
		t.Error("Expected no user in context")
	}
}

func TestGetUserFromContext_WithUser(t *testing.T) {
	claims := &services.Claims{
		UserID: 1,
		Login:  "testuser",
	}

	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), UserContextKey, claims)
	req = req.WithContext(ctx)

	user, ok := GetUserFromContext(req)
	if !ok {
		t.Error("Expected user in context")
	}

	if user.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", user.UserID)
	}

	if user.Login != "testuser" {
		t.Errorf("Expected login testuser, got %s", user.Login)
	}
}

func TestAuthMiddleware_WrongBearerFormat(t *testing.T) {
	mockDB := &MockDatabase{}
	authService := services.NewAuthService(mockDB, "test-secret")
	middleware := AuthMiddleware(authService)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	protectedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer")
	rr := httptest.NewRecorder()
	protectedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.Header.Set("Authorization", "Basic token")
	rr2 := httptest.NewRecorder()
	protectedHandler.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rr2.Code)
	}
}
