package services

import (
	"calculator/models"
	"fmt"
	"testing"
)

type MockDatabaseService struct {
	users  map[string]*models.User
	nextID int
}

func NewMockDatabaseService() *MockDatabaseService {
	return &MockDatabaseService{
		users:  make(map[string]*models.User),
		nextID: 1,
	}
}

func (m *MockDatabaseService) CreateUser(login, passwordHash string) (*models.User, error) {
	if _, exists := m.users[login]; exists {
		return nil, fmt.Errorf("user already exists")
	}

	user := &models.User{
		ID:           m.nextID,
		Login:        login,
		PasswordHash: passwordHash,
	}
	m.nextID++
	m.users[login] = user
	return user, nil
}

func (m *MockDatabaseService) GetUserByLogin(login string) (*models.User, error) {
	if user, exists := m.users[login]; exists {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockDatabaseService) Close() error {
	return nil
}

func TestAuthService_Register(t *testing.T) {
	db := NewMockDatabaseService()
	authService := NewAuthService(db, "test-secret")

	req := &models.RegisterRequest{
		Login:    "testuser",
		Password: "testpass123",
	}

	err := authService.Register(req)
	if err != nil {
		t.Errorf("Expected successful registration, got error: %v", err)
	}

	err = authService.Register(req)
	if err == nil {
		t.Error("Expected error for duplicate login, got nil")
	}

	shortLoginReq := &models.RegisterRequest{
		Login:    "ab",
		Password: "testpass123",
	}

	err = authService.Register(shortLoginReq)
	if err == nil {
		t.Error("Expected error for short login, got nil")
	}

	shortPassReq := &models.RegisterRequest{
		Login:    "testuser2",
		Password: "123",
	}

	err = authService.Register(shortPassReq)
	if err == nil {
		t.Error("Expected error for short password, got nil")
	}
}

func TestAuthService_Login(t *testing.T) {
	db := NewMockDatabaseService()
	authService := NewAuthService(db, "test-secret")

	registerReq := &models.RegisterRequest{
		Login:    "testuser",
		Password: "testpass123",
	}

	err := authService.Register(registerReq)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	loginReq := &models.LoginRequest{
		Login:    "testuser",
		Password: "testpass123",
	}

	response, err := authService.Login(loginReq)
	if err != nil {
		t.Errorf("Expected successful login, got error: %v", err)
	}

	if response == nil || response.Token == "" {
		t.Error("Expected token in response, got empty")
	}

	wrongPassReq := &models.LoginRequest{
		Login:    "testuser",
		Password: "wrongpass",
	}

	_, err = authService.Login(wrongPassReq)
	if err == nil {
		t.Error("Expected error for wrong password, got nil")
	}

	nonExistentReq := &models.LoginRequest{
		Login:    "nonexistent",
		Password: "testpass123",
	}

	_, err = authService.Login(nonExistentReq)
	if err == nil {
		t.Error("Expected error for non-existent user, got nil")
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	db := NewMockDatabaseService()
	authService := NewAuthService(db, "test-secret")

	registerReq := &models.RegisterRequest{
		Login:    "testuser",
		Password: "testpass123",
	}

	err := authService.Register(registerReq)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	loginReq := &models.LoginRequest{
		Login:    "testuser",
		Password: "testpass123",
	}

	response, err := authService.Login(loginReq)
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	claims, err := authService.ValidateToken(response.Token)
	if err != nil {
		t.Errorf("Expected valid token, got error: %v", err)
	}

	if claims == nil || claims.Login != "testuser" {
		t.Error("Expected valid claims with correct login")
	}

	_, err = authService.ValidateToken("invalid-token")
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}
