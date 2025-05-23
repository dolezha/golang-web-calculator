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

func getServerURL() string {
	if url := os.Getenv("SERVER_URL"); url != "" {
		return url
	}
	return "http://calc-service:8080"
}

var computingPower = getEnvInt("COMPUTING_POWER", 4)
var serverURL = getServerURL()

func StartAgent() {
	fmt.Printf("Starting agent with %d workers\n", computingPower)
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
		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		result := compute(task)
		err = submitTaskResult(task.ID, result)
		if err != nil {
			fmt.Printf("Error submitting task %s: %v\n", task.ID, err)
		}
	}
}

func getTaskFromServer() (*models.Task, error) {
	resp, err := http.Get(serverURL + "/internal/task")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no tasks")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var task models.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

func submitTaskResult(taskID string, result float64) error {
	url := fmt.Sprintf("%s/internal/task/%s", serverURL, taskID)
	payload := fmt.Sprintf(`{"result":%v}`, result)
	resp, err := http.Post(url, "application/json", strings.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error submitting result: %s", string(bodyBytes))
	}
	return nil
}

func compute(task *models.Task) float64 {
	getArgValue := func(arg string) float64 {
		if strings.HasPrefix(arg, "$") {
			taskID := strings.TrimPrefix(arg, "$")
			for {
				prevTask, err := getTaskResult(taskID)
				if err == nil && prevTask.Result != nil {
					return *prevTask.Result
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
		val, err := strconv.ParseFloat(arg, 64)
		if err != nil {
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
	resp, err := http.Get(fmt.Sprintf("%s/internal/task/%s", serverURL, taskID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("task not found")
	}

	var task models.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}
	return &task, nil
}
