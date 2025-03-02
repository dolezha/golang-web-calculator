package handlers

import (
	"calculator/services"
	"calculator/utils"
	"encoding/json"
	"net/http"
	"strings"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) > 2 && parts[len(parts)-2] == "task" {
			taskID := parts[len(parts)-1]
			task, exists := services.GetTask(taskID)
			if !exists {
				http.Error(w, "нет такой задачи", http.StatusNotFound)
				return
			}
			response := map[string]interface{}{
				"task": task,
			}
			utils.RespondWithJSON(w, response, http.StatusOK)
			return
		}

		task, err := services.GetNextTask()
		if err != nil {
			http.Error(w, "нет задач", http.StatusNotFound)
			return
		}
		response := map[string]interface{}{
			"task": task,
		}
		utils.RespondWithJSON(w, response, http.StatusOK)

	case http.MethodPost:
		var req struct {
			ID     string   `json:"id"`
			Result *float64 `json:"result"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Result == nil {
			utils.RespondWithJSON(w, map[string]interface{}{"error": "невалидные данные"}, http.StatusUnprocessableEntity)
			return
		}
		err := services.SubmitTaskResult(req.ID, *req.Result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		utils.RespondWithJSON(w, map[string]interface{}{"message": "результат записан"}, http.StatusOK)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
