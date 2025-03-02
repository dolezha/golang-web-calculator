package services

import (
	"calculator/models"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Добавляем новые типы для разбора выражений
type Operation struct {
	Type     string // "+", "-", "*", "/"
	Priority int    // приоритет операции
	Left     *Operation
	Right    *Operation
	Value    float64
	IsValue  bool
}

// Для упрощения используем глобальные карты. В реальной системе – БД.
var (
	expressions = make(map[string]*models.Expression)
	tasks       = make(map[string]*models.Task)
	mu          sync.Mutex
)

func CreateExpression(expr string) (*models.Expression, error) {
	if _, err := Calc(expr); err != nil {
		return nil, err
	}

	id := strconv.FormatInt(time.Now().UnixNano(), 10)
	expression := &models.Expression{
		ID:         id,
		Expression: expr,
		Status:     models.StatusPending,
	}
	mu.Lock()
	expressions[id] = expression
	mu.Unlock()

	err := splitExpressionIntoTasks(expression)
	if err != nil {
		return nil, err
	}
	return expression, nil
}

func splitExpressionIntoTasks(exp *models.Expression) error {
	fmt.Printf("Разбираем выражение: %s\n", exp.Expression)

	tree, err := parseExpression(exp.Expression)
	if err != nil {
		fmt.Printf("Ошибка разбора выражения: %v\n", err)
		return err
	}

	taskCounter := 1
	var createTasks func(*Operation) (string, error)
	createTasks = func(op *Operation) (string, error) {
		if op.IsValue {
			fmt.Printf("Создаем значение: %v\n", op.Value)
			return fmt.Sprintf("%v", op.Value), nil
		}

		leftArg, err := createTasks(op.Left)
		if err != nil {
			return "", err
		}

		rightArg, err := createTasks(op.Right)
		if err != nil {
			return "", err
		}

		taskID := fmt.Sprintf("%s_task%d", exp.ID, taskCounter)
		taskCounter++

		task := &models.Task{
			ID:            taskID,
			ExpressionID:  exp.ID,
			Arg1:          leftArg,
			Arg2:          rightArg,
			Operation:     op.Type,
			OperationTime: getOperationTime(op.Type),
			Status:        "pending",
		}

		fmt.Printf("Создаем задачу: %+v\n", task)

		mu.Lock()
		tasks[taskID] = task
		mu.Unlock()

		return fmt.Sprintf("$%s", taskID), nil
	}

	_, err = createTasks(tree)
	if err != nil {
		fmt.Printf("Ошибка создания задач: %v\n", err)
	}
	return err
}

func getOperationTime(op string) int64 {
	switch op {
	case "+":
		return getEnvInt64("TIME_ADDITION_MS", 1000)
	case "-":
		return getEnvInt64("TIME_SUBTRACTION_MS", 1000)
	case "*":
		return getEnvInt64("TIME_MULTIPLICATION_MS", 2000)
	case "/":
		return getEnvInt64("TIME_DIVISION_MS", 2000)
	default:
		return 1000
	}
}

func getEnvInt64(key string, fallback int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return i
		}
	}
	return fallback
}

func GetNextTask() (*models.Task, error) {
	mu.Lock()
	defer mu.Unlock()
	for _, task := range tasks {
		if task.Status == "pending" {
			task.Status = "computing"

			// Обновляем статус выражения
			if expr, exists := expressions[task.ExpressionID]; exists {
				expr.Status = models.StatusComputing
			}

			return task, nil
		}
	}
	return nil, errors.New("нет задач")
}

func SubmitTaskResult(taskID string, result float64) error {
	mu.Lock()
	defer mu.Unlock()

	fmt.Printf("Получен результат для задачи %s: %v\n", taskID, result)

	task, exists := tasks[taskID]
	if !exists {
		return errors.New("нет такой задачи")
	}
	if task.Status != "computing" {
		return errors.New("неверный статус задачи")
	}

	// Обновляем результат задачи
	task.Result = &result
	task.Status = "done"
	fmt.Printf("Задача %s обновлена: статус=%s, результат=%v\n",
		task.ID, task.Status, *task.Result)

	// Обновляем статус и результат выражения
	expr := expressions[task.ExpressionID]
	if expr != nil {
		fmt.Printf("Проверяем статус выражения %s\n", expr.ID)

		// Важное изменение: копируем задачи в локальный слайс
		var exprTasks []models.Task
		for _, t := range tasks {
			if t.ExpressionID == expr.ID {
				exprTasks = append(exprTasks, *t)
			}
		}

		// Проверяем все задачи выражения
		allDone := true
		var lastResult float64
		for _, t := range exprTasks {
			fmt.Printf("  Задача %s: статус=%s\n", t.ID, t.Status)
			if t.Status != "done" {
				allDone = false
				break
			}
			if t.Result != nil {
				lastResult = *t.Result
			}
		}

		// Если все задачи выполнены, обновляем выражение
		if allDone {
			expr.Status = models.StatusDone
			expr.Result = &lastResult
			fmt.Printf("Выражение %s завершено с результатом %v\n",
				expr.ID, *expr.Result)
		} else {
			expr.Status = models.StatusComputing
			fmt.Printf("Выражение %s в процессе вычисления\n", expr.ID)
		}
	}

	return nil
}

func GetExpression(id string) (*models.Expression, bool) {
	mu.Lock()
	defer mu.Unlock()
	exp, exists := expressions[id]
	return exp, exists
}

func GetExpressionsList() []*models.Expression {
	mu.Lock()
	defer mu.Unlock()
	list := []*models.Expression{}
	for _, exp := range expressions {
		list = append(list, exp)
	}
	return list
}

func GetTask(id string) (*models.Task, bool) {
	mu.Lock()
	defer mu.Unlock()
	task, exists := tasks[id]
	return task, exists
}

func buildTree(tokens []string) (*Operation, error) {
	if len(tokens) == 0 {
		return nil, errors.New("пустое выражение")
	}

	// Сначала ищем + и -
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i] == "+" || tokens[i] == "-" {
			left, err := buildTree(tokens[:i])
			if err != nil {
				return nil, fmt.Errorf("ошибка в левой части: %v", err)
			}
			right, err := buildTree(tokens[i+1:])
			if err != nil {
				return nil, fmt.Errorf("ошибка в правой части: %v", err)
			}
			return &Operation{
				Type:     tokens[i],
				Priority: 1,
				Left:     left,
				Right:    right,
			}, nil
		}
	}

	// Затем * и /
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i] == "*" || tokens[i] == "/" {
			left, err := buildTree(tokens[:i])
			if err != nil {
				return nil, fmt.Errorf("ошибка в левой части: %v", err)
			}
			right, err := buildTree(tokens[i+1:])
			if err != nil {
				return nil, fmt.Errorf("ошибка в правой части: %v", err)
			}
			return &Operation{
				Type:     tokens[i],
				Priority: 2,
				Left:     left,
				Right:    right,
			}, nil
		}
	}

	// Если нет операторов, это число
	val, err := strconv.ParseFloat(strings.Join(tokens, ""), 64)
	if err != nil {
		return nil, fmt.Errorf("некорректное число: %v", err)
	}
	return &Operation{
		IsValue: true,
		Value:   val,
	}, nil
}

func parseExpression(expr string) (*Operation, error) {
	tokens := tokenize(expr)
	return buildTree(tokens)
}

func tokenize(expr string) []string {
	var tokens []string
	var num strings.Builder

	for _, c := range expr {
		switch c {
		case '+', '-', '*', '/':
			if num.Len() > 0 {
				tokens = append(tokens, num.String())
				num.Reset()
			}
			tokens = append(tokens, string(c))
		case ' ':
			continue
		default:
			num.WriteRune(c)
		}
	}
	if num.Len() > 0 {
		tokens = append(tokens, num.String())
	}
	return tokens
}
