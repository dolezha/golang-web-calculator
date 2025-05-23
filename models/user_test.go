package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUser_JSONSerialization(t *testing.T) {
	user := User{
		ID:           1,
		Login:        "testuser",
		PasswordHash: "secret",
		CreatedAt:    time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}

	jsonStr := string(data)
	if contains(jsonStr, "secret") {
		t.Error("Password hash should not be serialized to JSON")
	}

	if !contains(jsonStr, "testuser") {
		t.Error("Login should be serialized to JSON")
	}

	var unmarshaled User
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal user: %v", err)
	}

	if unmarshaled.Login != user.Login {
		t.Errorf("Expected login %s, got %s", user.Login, unmarshaled.Login)
	}
}

func TestRegisterRequest_Validation(t *testing.T) {
	req := RegisterRequest{
		Login:    "testuser",
		Password: "password123",
	}

	if req.Login == "" {
		t.Error("Login should not be empty")
	}

	if req.Password == "" {
		t.Error("Password should not be empty")
	}
}

func TestLoginRequest_Validation(t *testing.T) {
	req := LoginRequest{
		Login:    "testuser",
		Password: "password123",
	}

	if req.Login == "" {
		t.Error("Login should not be empty")
	}

	if req.Password == "" {
		t.Error("Password should not be empty")
	}
}

func TestLoginResponse_JSONSerialization(t *testing.T) {
	resp := LoginResponse{
		Token: "jwt-token-here",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	if !contains(string(data), "jwt-token-here") {
		t.Error("Token should be serialized to JSON")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsInner(s, substr)))
}

func containsInner(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
