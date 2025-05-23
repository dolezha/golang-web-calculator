package services

import (
	"calculator/models"
	"fmt"
	"strconv"
	"time"
)

type ExpressionService struct {
	db *DatabaseService
}

func NewExpressionService(db *DatabaseService) *ExpressionService {
	return &ExpressionService{db: db}
}

func (es *ExpressionService) CreateExpression(userID int, expr string) (*models.Expression, error) {
	if _, err := Calc(expr); err != nil {
		return nil, fmt.Errorf("invalid expression: %v", err)
	}

	id := strconv.FormatInt(time.Now().UnixNano(), 10)
	expression := &models.Expression{
		ID:         id,
		UserID:     userID,
		Expression: expr,
		Status:     models.StatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := es.db.CreateExpression(expression); err != nil {
		return nil, fmt.Errorf("error saving expression: %v", err)
	}

	if err := es.splitExpressionIntoTasks(expression); err != nil {
		return nil, fmt.Errorf("error creating tasks: %v", err)
	}

	return expression, nil
}

func (es *ExpressionService) GetExpression(id string, userID int) (*models.Expression, error) {
	return es.db.GetExpression(id, userID)
}

func (es *ExpressionService) GetUserExpressions(userID int) ([]*models.Expression, error) {
	return es.db.GetUserExpressions(userID)
}

func (es *ExpressionService) splitExpressionIntoTasks(exp *models.Expression) error {
	tree, err := parseExpression(exp.Expression)
	if err != nil {
		return err
	}

	taskCounter := 1
	var createTasks func(*Operation) (string, error)
	createTasks = func(op *Operation) (string, error) {
		if op.IsValue {
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
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := es.db.CreateTask(task); err != nil {
			return "", fmt.Errorf("error saving task: %v", err)
		}

		return fmt.Sprintf("$%s", taskID), nil
	}

	_, err = createTasks(tree)
	return err
}

func (es *ExpressionService) GetNextTask() (*models.Task, error) {
	tasks, err := es.db.GetPendingTasks()
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, fmt.Errorf("no available tasks")
	}

	task := tasks[0]
	task.Status = "computing"
	task.UpdatedAt = time.Now()

	if err := es.db.UpdateTask(task); err != nil {
		return nil, fmt.Errorf("error updating task: %v", err)
	}

	expr, err := es.db.GetExpression(task.ExpressionID, 0)
	if err == nil {
		expr.Status = models.StatusComputing
		es.db.UpdateExpression(expr)
	}

	return task, nil
}

func (es *ExpressionService) GetTaskByID(taskID string) (*models.Task, error) {
	return es.db.GetTask(taskID)
}

func (es *ExpressionService) SubmitTaskResult(taskID string, result float64) error {
	task, err := es.db.GetTask(taskID)
	if err != nil {
		return fmt.Errorf("task not found: %v", err)
	}

	if task.Status != "computing" {
		return fmt.Errorf("invalid task status: %s", task.Status)
	}

	task.Result = &result
	task.Status = "done"
	task.UpdatedAt = time.Now()

	if err := es.db.UpdateTask(task); err != nil {
		return fmt.Errorf("error updating task: %v", err)
	}

	return es.checkExpressionCompletion(task.ExpressionID)
}

func (es *ExpressionService) checkExpressionCompletion(expressionID string) error {
	expr, err := es.db.GetExpression(expressionID, 0)
	if err != nil {
		return fmt.Errorf("expression not found: %v", err)
	}

	exprTasks, err := es.db.GetTasksByExpressionID(expressionID)
	if err != nil {
		return fmt.Errorf("error getting tasks: %v", err)
	}

	allDone := true
	var lastResult float64
	for _, task := range exprTasks {
		if task.Status != "done" {
			allDone = false
			break
		}
		if task.Result != nil {
			lastResult = *task.Result
		}
	}

	if allDone {
		expr.Status = models.StatusDone
		expr.Result = &lastResult
	} else {
		expr.Status = models.StatusComputing
	}

	expr.UpdatedAt = time.Now()
	return es.db.UpdateExpression(expr)
}
