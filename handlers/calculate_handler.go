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
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
