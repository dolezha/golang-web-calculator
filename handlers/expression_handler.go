package handlers

import (
	"calculator/models"
	"calculator/services"
	"calculator/utils"
	"net/http"
	"strings"
)

func ExpressionsListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	expressions := services.GetExpressionsList()
	response := struct {
		Expressions []struct {
			ID     string                  `json:"id"`
			Status models.ExpressionStatus `json:"status"`
			Result *float64                `json:"result,omitempty"`
		} `json:"expressions"`
	}{
		Expressions: make([]struct {
			ID     string                  `json:"id"`
			Status models.ExpressionStatus `json:"status"`
			Result *float64                `json:"result,omitempty"`
		}, len(expressions)),
	}

	for i, exp := range expressions {
		response.Expressions[i].ID = exp.ID
		response.Expressions[i].Status = exp.Status
		response.Expressions[i].Result = exp.Result
	}

	utils.RespondWithJSON(w, response, http.StatusOK)
}

func ExpressionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	id := parts[len(parts)-1]
	exp, exists := services.GetExpression(id)
	if !exists {
		http.Error(w, "нет такого выражения", http.StatusNotFound)
		return
	}
	response := struct {
		Expression struct {
			ID     string                  `json:"id"`
			Status models.ExpressionStatus `json:"status"`
			Result *float64                `json:"result,omitempty"`
		} `json:"expression"`
	}{
		Expression: struct {
			ID     string                  `json:"id"`
			Status models.ExpressionStatus `json:"status"`
			Result *float64                `json:"result,omitempty"`
		}{
			ID:     exp.ID,
			Status: exp.Status,
			Result: exp.Result,
		},
	}
	utils.RespondWithJSON(w, response, http.StatusOK)
}
