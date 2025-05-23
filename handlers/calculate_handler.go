package handlers

import (
	"calculator/middleware"
	"calculator/models"
	"calculator/services"
	"calculator/utils"
	"encoding/json"
	"net/http"
)

type CalculateHandler struct {
	expressionService *services.ExpressionService
}

func NewCalculateHandler(expressionService *services.ExpressionService) *CalculateHandler {
	return &CalculateHandler{expressionService: expressionService}
}

func (ch *CalculateHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		utils.RespondWithJSON(w, map[string]string{"error": "User not authorized"}, http.StatusUnauthorized)
		return
	}

	var reqBody models.RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		utils.RespondWithJSON(w, map[string]string{"error": "Invalid request body"}, http.StatusUnprocessableEntity)
		return
	}

	expression, err := ch.expressionService.CreateExpression(claims.UserID, reqBody.Expression)
	if err != nil {
		utils.RespondWithJSON(w, map[string]string{"error": err.Error()}, http.StatusUnprocessableEntity)
		return
	}

	response := map[string]string{
		"id": expression.ID,
	}

	if err := utils.RespondWithJSON(w, response, http.StatusCreated); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
