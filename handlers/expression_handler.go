package handlers

import (
	"calculator/middleware"
	"calculator/services"
	"calculator/utils"
	"net/http"
	"strings"
)

type ExpressionHandler struct {
	expressionService *services.ExpressionService
}

func NewExpressionHandler(expressionService *services.ExpressionService) *ExpressionHandler {
	return &ExpressionHandler{expressionService: expressionService}
}

func (eh *ExpressionHandler) GetExpression(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		utils.RespondWithJSON(w, map[string]string{"error": "Пользователь не авторизован"}, http.StatusUnauthorized)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")
	if path == "" {
		utils.RespondWithJSON(w, map[string]string{"error": "ID выражения не указан"}, http.StatusBadRequest)
		return
	}

	expression, err := eh.expressionService.GetExpression(path, claims.UserID)
	if err != nil {
		utils.RespondWithJSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
		return
	}

	utils.RespondWithJSON(w, expression, http.StatusOK)
}

func (eh *ExpressionHandler) GetExpressions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		utils.RespondWithJSON(w, map[string]string{"error": "Пользователь не авторизован"}, http.StatusUnauthorized)
		return
	}

	expressions, err := eh.expressionService.GetUserExpressions(claims.UserID)
	if err != nil {
		utils.RespondWithJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, expressions, http.StatusOK)
}
