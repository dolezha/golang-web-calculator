package main

import (
	"bytes"
	"calculator/handlers"
	"calculator/middleware"
	"calculator/services"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestFullUserWorkflow(t *testing.T) {
	dbPath := "./test_integration.db"
	defer os.Remove(dbPath)

	db, err := services.NewDatabaseService(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	authService := services.NewAuthService(db, "test-secret-key")
	expressionService := services.NewExpressionService(db)

	authHandler := handlers.NewAuthHandler(authService)
	calculateHandler := handlers.NewCalculateHandler(expressionService)
	expressionHandler := handlers.NewExpressionHandler(expressionService)

	authMiddleware := middleware.AuthMiddleware(authService)

	t.Run("User Registration", func(t *testing.T) {
		reqBody := map[string]string{
			"login":    "testuser",
			"password": "testpass123",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		authHandler.Register(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d. Body: %s", rr.Code, rr.Body.String())
		}
	})

	var jwtToken string
	t.Run("User Login", func(t *testing.T) {
		reqBody := map[string]string{
			"login":    "testuser",
			"password": "testpass123",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		authHandler.Login(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}

		var response map[string]string
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse login response: %v", err)
		}

		jwtToken = response["token"]
		if jwtToken == "" {
			t.Fatal("JWT token not received")
		}
	})

	var expressionID string
	t.Run("Create Expression", func(t *testing.T) {
		reqBody := map[string]string{
			"expression": "2+3*4",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+jwtToken)

		handler := authMiddleware(http.HandlerFunc(calculateHandler.Calculate))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf("Expected status 201, got %d. Body: %s", rr.Code, rr.Body.String())
		}

		var response map[string]string
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse calculate response: %v", err)
		}

		expressionID = response["id"]
		if expressionID == "" {
			t.Fatal("Expression ID not received")
		}
	})

	t.Run("Get Expression Result", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)

		req := httptest.NewRequest("GET", "/api/v1/expressions/"+expressionID, nil)
		req.Header.Set("Authorization", "Bearer "+jwtToken)

		handler := authMiddleware(http.HandlerFunc(expressionHandler.GetExpression))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}

		var response map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse expression response: %v", err)
		}

		if response["id"] != expressionID {
			t.Errorf("Expected expression ID %s, got %v", expressionID, response["id"])
		}

		if response["expression"] != "2+3*4" {
			t.Errorf("Expected expression '2+3*4', got %v", response["expression"])
		}
	})

	t.Run("Get User Expressions", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/expressions", nil)
		req.Header.Set("Authorization", "Bearer "+jwtToken)

		handler := authMiddleware(http.HandlerFunc(expressionHandler.GetExpressions))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", rr.Code, rr.Body.String())
		}

		var expressions []map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &expressions); err != nil {
			t.Fatalf("Failed to parse expressions response: %v", err)
		}

		if len(expressions) != 1 {
			t.Errorf("Expected 1 expression, got %d", len(expressions))
		}

		if len(expressions) > 0 && expressions[0]["id"] != expressionID {
			t.Errorf("Expected expression ID %s, got %v", expressionID, expressions[0]["id"])
		}
	})
}

func TestUserIsolation(t *testing.T) {
	dbPath := "./test_isolation.db"
	defer os.Remove(dbPath)

	db, err := services.NewDatabaseService(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	authService := services.NewAuthService(db, "test-secret-key")
	expressionService := services.NewExpressionService(db)

	authHandler := handlers.NewAuthHandler(authService)
	calculateHandler := handlers.NewCalculateHandler(expressionService)
	expressionHandler := handlers.NewExpressionHandler(expressionService)

	authMiddleware := middleware.AuthMiddleware(authService)

	users := []struct {
		login    string
		password string
		token    string
	}{
		{"user1", "pass1", ""},
		{"user2", "pass2", ""},
	}

	for i := range users {
		reqBody := map[string]string{
			"login":    users[i].login,
			"password": users[i].password,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		authHandler.Register(rr, req)
		if rr.Code != http.StatusCreated {
			t.Fatalf("Failed to register user %s: %d", users[i].login, rr.Code)
		}

		req = httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()

		authHandler.Login(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("Failed to login user %s: %d", users[i].login, rr.Code)
		}

		var response map[string]string
		json.Unmarshal(rr.Body.Bytes(), &response)
		users[i].token = response["token"]
	}

	for i, user := range users {
		reqBody := map[string]string{
			"expression": "1+" + string(rune('0'+i+1)),
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+user.token)

		handler := authMiddleware(http.HandlerFunc(calculateHandler.Calculate))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Fatalf("Failed to create expression for user %s: %d", user.login, rr.Code)
		}
	}

	for i, user := range users {
		req := httptest.NewRequest("GET", "/api/v1/expressions", nil)
		req.Header.Set("Authorization", "Bearer "+user.token)

		handler := authMiddleware(http.HandlerFunc(expressionHandler.GetExpressions))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("Failed to get expressions for user %s: %d", user.login, rr.Code)
		}

		var expressions []map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &expressions)

		if len(expressions) != 1 {
			t.Errorf("User %s should see exactly 1 expression, got %d", user.login, len(expressions))
		}

		if len(expressions) > 0 {
			expectedExpr := "1+" + string(rune('0'+i+1))
			if expressions[0]["expression"] != expectedExpr {
				t.Errorf("User %s should see expression '%s', got '%v'", user.login, expectedExpr, expressions[0]["expression"])
			}
		}
	}
}

func TestErrorHandling(t *testing.T) {
	dbPath := "./test_errors.db"
	defer os.Remove(dbPath)

	db, err := services.NewDatabaseService(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	authService := services.NewAuthService(db, "test-secret-key")
	authHandler := handlers.NewAuthHandler(authService)

	t.Run("Duplicate User Registration", func(t *testing.T) {
		reqBody := map[string]string{
			"login":    "duplicate",
			"password": "pass123",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		authHandler.Register(rr, req)
		if rr.Code != http.StatusCreated {
			t.Fatalf("First registration should succeed: %d", rr.Code)
		}

		req = httptest.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()

		authHandler.Register(rr, req)
		if rr.Code == http.StatusCreated {
			t.Error("Duplicate registration should fail")
		}
	})

	t.Run("Invalid Login Credentials", func(t *testing.T) {
		reqBody := map[string]string{
			"login":    "nonexistent",
			"password": "wrongpass",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		authHandler.Login(rr, req)
		if rr.Code == http.StatusOK {
			t.Error("Login with invalid credentials should fail")
		}
	})
}
