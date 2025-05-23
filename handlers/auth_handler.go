package handlers

import (
	"calculator/models"
	"calculator/services"
	"calculator/utils"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithJSON(w, map[string]string{"error": "Invalid request format"}, http.StatusBadRequest)
		return
	}

	if err := ah.authService.Register(&req); err != nil {
		utils.RespondWithJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	utils.RespondWithJSON(w, map[string]string{"message": "User registered successfully"}, http.StatusCreated)
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithJSON(w, map[string]string{"error": "Invalid request format"}, http.StatusBadRequest)
		return
	}

	response, err := ah.authService.Login(&req)
	if err != nil {
		utils.RespondWithJSON(w, map[string]string{"error": err.Error()}, http.StatusUnauthorized)
		return
	}

	utils.RespondWithJSON(w, response, http.StatusOK)
}
