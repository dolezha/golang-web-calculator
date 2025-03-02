package agent

import (
	"calculator/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

var computingPower = getEnvInt("COMPUTING_POWER", 4)

func StartAgent() {
	for i := 0; i < computingPower; i++ {
		go worker(i)
	}
	select {}
}

func worker(id int) {
	for {
		task, err := getTaskFromServer()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Printf("Воркер %d получил задачу %s\n", id, task.ID)
		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		result := compute(task)
		err = submitTaskResult(task.ID, result)
		if err != nil {
			fmt.Printf("Ошибка отправки результата задачи %s: %v\n", task.ID, err)
		} else {
			fmt.Printf("Воркер %d завершил задачу %s с результатом %v\n", id, task.ID, result)
		}
	}
}

func getTaskFromServer() (*models.Task, error) {
	resp, err := http.Get("http://localhost:8080/internal/task")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("нет задач")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var data struct {
		Task models.Task `json:"task"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data.Task, nil
}

func submitTaskResult(taskID string, result float64) error {
	url := "http://localhost:8080/internal/task"
	payload := fmt.Sprintf(`{"id":"%s", "result":%v}`, taskID, result)
	resp, err := http.Post(url, "application/json", strings.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка отправки результата: %s", string(bodyBytes))
	}
	return nil
}

func compute(task *models.Task) float64 {
	fmt.Printf("Вычисляю задачу: %+v\n", task)

	getArgValue := func(arg string) float64 {
		if strings.HasPrefix(arg, "$") {
			taskID := strings.TrimPrefix(arg, "$")
			fmt.Printf("Получаем результат задачи %s\n", taskID)
			for {
				prevTask, err := getTaskResult(taskID)
				if err == nil && prevTask.Result != nil {
					fmt.Printf("Получен результат задачи %s: %v\n", taskID, *prevTask.Result)
					return *prevTask.Result
				}
				fmt.Printf("Ожидание результата задачи %s...\n", taskID)
				time.Sleep(100 * time.Millisecond)
			}
		}
		val, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			fmt.Printf("Ошибка преобразования аргумента %s: %v\n", arg, err)
			return 0
		}
		return val
	}

	switch task.Operation {
	case "+":
		return getArgValue(task.Arg1) + getArgValue(task.Arg2)
	case "-":
		return getArgValue(task.Arg1) - getArgValue(task.Arg2)
	case "*":
		return getArgValue(task.Arg1) * getArgValue(task.Arg2)
	case "/":
		b := getArgValue(task.Arg2)
		if b == 0 {
			return 0
		}
		return getArgValue(task.Arg1) / b
	default:
		return 0
	}
}

func getTaskResult(taskID string) (*models.Task, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/internal/task/%s", taskID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		Task models.Task `json:"task"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data.Task, nil
}
