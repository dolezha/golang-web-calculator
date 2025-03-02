package handlers

import (
	"calculator/models"
	"calculator/services"
	"calculator/utils"
	"encoding/json"
	"net/http"
)

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody models.RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		utils.RespondWithJSON(w, map[string]string{"error": "Invalid request body"}, http.StatusUnprocessableEntity)
		return
	}

	expression, err := services.CreateExpression(reqBody.Expression)
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
