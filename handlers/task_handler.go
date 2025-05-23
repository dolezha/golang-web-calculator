package handlers

import (
	"calculator/services"
	"calculator/utils"
	"encoding/json"
	"net/http"
	"strings"
)

type TaskHandler struct {
	expressionService *services.ExpressionService
}

func NewTaskHandler(expressionService *services.ExpressionService) *TaskHandler {
	return &TaskHandler{expressionService: expressionService}
}

func (th *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	task, err := th.expressionService.GetNextTask()
	if err != nil {
		utils.RespondWithJSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
		return
	}

	utils.RespondWithJSON(w, task, http.StatusOK)
}

func (th *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/internal/task/")
	if path == "" {
		utils.RespondWithJSON(w, map[string]string{"error": "ID задачи не указан"}, http.StatusBadRequest)
		return
	}

	task, err := th.expressionService.GetTaskByID(path)
	if err != nil {
		utils.RespondWithJSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
		return
	}

	utils.RespondWithJSON(w, task, http.StatusOK)
}

func (th *TaskHandler) SubmitTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/internal/task/")
	if path == "" {
		utils.RespondWithJSON(w, map[string]string{"error": "ID задачи не указан"}, http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Result float64 `json:"result"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		utils.RespondWithJSON(w, map[string]string{"error": "Неверный формат запроса"}, http.StatusBadRequest)
		return
	}

	if err := th.expressionService.SubmitTaskResult(path, reqBody.Result); err != nil {
		utils.RespondWithJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	utils.RespondWithJSON(w, map[string]string{"message": "Результат принят"}, http.StatusOK)
}
