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
<<<<<<< HEAD
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
=======
		response := models.ResponseBody{
			Error: utils.ToPtr("Invalid request body"),
		}
		utils.RespondWithJSON(w, response, http.StatusUnprocessableEntity)
		return
	}

	result, err := services.Calc(reqBody.Expression)
	if err != nil {
		response := models.ResponseBody{
			Error: utils.ToPtr("Expression is not valid"),
		}
		if err.Error() == "деление на ноль" {
			response.Error = utils.ToPtr("Division by zero")
		}
		utils.RespondWithJSON(w, response, http.StatusUnprocessableEntity)
		return
	}

	response := models.ResponseBody{
		Result: &result,
	}
	if err := utils.RespondWithJSON(w, response, http.StatusOK); err != nil {
>>>>>>> c030495b399869c8d46fd81d0464189aaf4466c3
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
