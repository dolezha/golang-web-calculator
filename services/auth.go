package services

import (
	"calculator/models"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type DatabaseInterface interface {
	CreateUser(login, passwordHash string) (*models.User, error)
	GetUserByLogin(login string) (*models.User, error)
	Close() error
}

type AuthService struct {
	db        DatabaseInterface
	jwtSecret []byte
}

type Claims struct {
	UserID int    `json:"user_id"`
	Login  string `json:"login"`
	jwt.RegisteredClaims
}

func NewAuthService(db DatabaseInterface, jwtSecret string) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: []byte(jwtSecret),
	}
}

func (as *AuthService) Register(req *models.RegisterRequest) error {
	_, err := as.db.GetUserByLogin(req.Login)
	if err == nil {
		return fmt.Errorf("user with this login already exists")
	}

	if len(req.Login) < 3 {
		return fmt.Errorf("login must contain at least 3 characters")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("password must contain at least 6 characters")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("password hashing error: %v", err)
	}

	_, err = as.db.CreateUser(req.Login, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("user creation error: %v", err)
	}

	return nil
}

func (as *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := as.db.GetUserByLogin(req.Login)
	if err != nil {
		return nil, fmt.Errorf("invalid login or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid login or password")
	}

	token, err := as.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("token creation error: %v", err)
	}

	return &models.LoginResponse{Token: token}, nil
}

func (as *AuthService) generateToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Login:  user.Login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(as.jwtSecret)
}

func (as *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return as.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
