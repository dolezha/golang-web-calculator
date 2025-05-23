package services

import (
	"calculator/models"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDatabaseService_CreateUser_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	service := &DatabaseService{db: db}

	mock.ExpectExec("INSERT INTO users").
		WithArgs("testuser", "hashedpassword").
		WillReturnResult(sqlmock.NewResult(1, 1))

	user, err := service.CreateUser("testuser", "hashedpassword")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.Login != "testuser" {
		t.Errorf("Expected login 'testuser', got '%s'", user.Login)
	}
	if user.ID != 1 {
		t.Errorf("Expected ID 1, got %d", user.ID)
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs("erroruser", "hash").
		WillReturnError(errors.New("database error"))

	_, err = service.CreateUser("erroruser", "hash")
	if err == nil {
		t.Error("Expected error for database failure")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDatabaseService_GetUserByLogin_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	service := &DatabaseService{db: db}

	rows := sqlmock.NewRows([]string{"id", "login", "password_hash", "created_at"}).
		AddRow(1, "testuser", "hashedpassword", time.Now())

	mock.ExpectQuery("SELECT id, login, password_hash, created_at FROM users WHERE login = ?").
		WithArgs("testuser").
		WillReturnRows(rows)

	user, err := service.GetUserByLogin("testuser")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if user.Login != "testuser" {
		t.Errorf("Expected login 'testuser', got '%s'", user.Login)
	}

	mock.ExpectQuery("SELECT id, login, password_hash, created_at FROM users WHERE login = ?").
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	_, err = service.GetUserByLogin("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existing user")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDatabaseService_CreateExpression_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	service := &DatabaseService{db: db}

	expr := &models.Expression{
		ID:         "test-id",
		UserID:     1,
		Expression: "2+2",
		Status:     models.StatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	mock.ExpectExec("INSERT INTO expressions").
		WithArgs(expr.ID, expr.UserID, expr.Expression, expr.Status, expr.CreatedAt, expr.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = service.CreateExpression(expr)
	if err != nil {
		t.Fatalf("Failed to create expression: %v", err)
	}

	mock.ExpectExec("INSERT INTO expressions").
		WithArgs(expr.ID, expr.UserID, expr.Expression, expr.Status, expr.CreatedAt, expr.UpdatedAt).
		WillReturnError(errors.New("database error"))

	err = service.CreateExpression(expr)
	if err == nil {
		t.Error("Expected error for database failure")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDatabaseService_GetExpression_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	service := &DatabaseService{db: db}

	rows := sqlmock.NewRows([]string{"id", "user_id", "expression", "status", "result", "created_at", "updated_at"}).
		AddRow("test-id", 1, "2+2", "pending", nil, time.Now(), time.Now())

	mock.ExpectQuery("SELECT id, user_id, expression, status, result, created_at, updated_at FROM expressions WHERE id = \\? AND user_id = \\?").
		WithArgs("test-id", 1).
		WillReturnRows(rows)

	expr, err := service.GetExpression("test-id", 1)
	if err != nil {
		t.Fatalf("Failed to get expression: %v", err)
	}

	if expr.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", expr.ID)
	}

	rows2 := sqlmock.NewRows([]string{"id", "user_id", "expression", "status", "result", "created_at", "updated_at"}).
		AddRow("test-id", 1, "2+2", "pending", nil, time.Now(), time.Now())

	mock.ExpectQuery("SELECT id, user_id, expression, status, result, created_at, updated_at FROM expressions WHERE id = \\?").
		WithArgs("test-id").
		WillReturnRows(rows2)

	expr2, err := service.GetExpression("test-id", 0)
	if err != nil {
		t.Fatalf("Failed to get expression without userID: %v", err)
	}

	if expr2.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", expr2.ID)
	}

	mock.ExpectQuery("SELECT id, user_id, expression, status, result, created_at, updated_at FROM expressions WHERE id = \\? AND user_id = \\?").
		WithArgs("nonexistent", 1).
		WillReturnError(sql.ErrNoRows)

	_, err = service.GetExpression("nonexistent", 1)
	if err == nil {
		t.Error("Expected error for non-existing expression")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDatabaseService_UpdateExpression_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	service := &DatabaseService{db: db}

	expr := &models.Expression{
		ID:     "test-id",
		Status: models.StatusComputing,
		Result: &[]float64{4.0}[0],
	}

	mock.ExpectExec("UPDATE expressions SET status = \\?, result = \\?, updated_at = \\? WHERE id = \\?").
		WithArgs(expr.Status, expr.Result, sqlmock.AnyArg(), expr.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = service.UpdateExpression(expr)
	if err != nil {
		t.Fatalf("Failed to update expression: %v", err)
	}

	mock.ExpectExec("UPDATE expressions SET status = \\?, result = \\?, updated_at = \\? WHERE id = \\?").
		WithArgs(expr.Status, expr.Result, sqlmock.AnyArg(), expr.ID).
		WillReturnError(errors.New("database error"))

	err = service.UpdateExpression(expr)
	if err == nil {
		t.Error("Expected error for database failure")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDatabaseService_GetUserExpressions_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	service := &DatabaseService{db: db}

	rows := sqlmock.NewRows([]string{"id", "user_id", "expression", "status", "result", "created_at", "updated_at"}).
		AddRow("test-id-1", 1, "2+2", "pending", nil, time.Now(), time.Now()).
		AddRow("test-id-2", 1, "3+3", "pending", nil, time.Now(), time.Now())

	mock.ExpectQuery("SELECT id, user_id, expression, status, result, created_at, updated_at FROM expressions WHERE user_id = \\? ORDER BY created_at DESC").
		WithArgs(1).
		WillReturnRows(rows)

	expressions, err := service.GetUserExpressions(1)
	if err != nil {
		t.Fatalf("Failed to get user expressions: %v", err)
	}

	if len(expressions) != 2 {
		t.Errorf("Expected 2 expressions, got %d", len(expressions))
	}

	mock.ExpectQuery("SELECT id, user_id, expression, status, result, created_at, updated_at FROM expressions WHERE user_id = \\? ORDER BY created_at DESC").
		WithArgs(1).
		WillReturnError(errors.New("database error"))

	_, err = service.GetUserExpressions(1)
	if err == nil {
		t.Error("Expected error for database failure")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDatabaseService_CreateTask_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	service := &DatabaseService{db: db}

	task := &models.Task{
		ID:            "task-id",
		ExpressionID:  "expr-id",
		Arg1:          "2",
		Arg2:          "2",
		Operation:     "+",
		OperationTime: 1000,
		Status:        "pending",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	mock.ExpectExec("INSERT INTO tasks").
		WithArgs(task.ID, task.ExpressionID, task.Arg1, task.Arg2, task.Operation, task.OperationTime, task.Status, task.CreatedAt, task.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = service.CreateTask(task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	mock.ExpectExec("INSERT INTO tasks").
		WithArgs(task.ID, task.ExpressionID, task.Arg1, task.Arg2, task.Operation, task.OperationTime, task.Status, task.CreatedAt, task.UpdatedAt).
		WillReturnError(errors.New("database error"))

	err = service.CreateTask(task)
	if err == nil {
		t.Error("Expected error for database failure")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDatabaseService_GetTask_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	service := &DatabaseService{db: db}

	rows := sqlmock.NewRows([]string{"id", "expression_id", "arg1", "arg2", "operation", "operation_time", "status", "result", "created_at", "updated_at"}).
		AddRow("task-id", "expr-id", "2", "2", "+", 1000, "pending", nil, time.Now(), time.Now())

	mock.ExpectQuery("SELECT id, expression_id, arg1, arg2, operation, operation_time, status, result, created_at, updated_at FROM tasks WHERE id = \\?").
		WithArgs("task-id").
		WillReturnRows(rows)

	task, err := service.GetTask("task-id")
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	if task.ID != "task-id" {
		t.Errorf("Expected ID 'task-id', got '%s'", task.ID)
	}

	mock.ExpectQuery("SELECT id, expression_id, arg1, arg2, operation, operation_time, status, result, created_at, updated_at FROM tasks WHERE id = \\?").
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	_, err = service.GetTask("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existing task")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDatabaseService_UpdateTask_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	service := &DatabaseService{db: db}

	task := &models.Task{
		ID:     "task-id",
		Status: "completed",
		Result: &[]float64{4.0}[0],
	}

	mock.ExpectExec("UPDATE tasks SET status = \\?, result = \\?, updated_at = \\? WHERE id = \\?").
		WithArgs(task.Status, task.Result, sqlmock.AnyArg(), task.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = service.UpdateTask(task)
	if err != nil {
		t.Fatalf("Failed to update task: %v", err)
	}

	mock.ExpectExec("UPDATE tasks SET status = \\?, result = \\?, updated_at = \\? WHERE id = \\?").
		WithArgs(task.Status, task.Result, sqlmock.AnyArg(), task.ID).
		WillReturnError(errors.New("database error"))

	err = service.UpdateTask(task)
	if err == nil {
		t.Error("Expected error for database failure")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDatabaseService_GetPendingTasks_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	service := &DatabaseService{db: db}

	rows := sqlmock.NewRows([]string{"id", "expression_id", "arg1", "arg2", "operation", "operation_time", "status", "result", "created_at", "updated_at"}).
		AddRow("task-id-1", "expr-id", "2", "2", "+", 1000, "pending", nil, time.Now(), time.Now())

	mock.ExpectQuery("SELECT id, expression_id, arg1, arg2, operation, operation_time, status, result, created_at, updated_at FROM tasks WHERE status = 'pending' ORDER BY created_at ASC").
		WillReturnRows(rows)

	tasks, err := service.GetPendingTasks()
	if err != nil {
		t.Fatalf("Failed to get pending tasks: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Expected 1 pending task, got %d", len(tasks))
	}

	if tasks[0].ID != "task-id-1" {
		t.Errorf("Expected task ID 'task-id-1', got '%s'", tasks[0].ID)
	}

	mock.ExpectQuery("SELECT id, expression_id, arg1, arg2, operation, operation_time, status, result, created_at, updated_at FROM tasks WHERE status = 'pending' ORDER BY created_at ASC").
		WillReturnError(errors.New("database error"))

	_, err = service.GetPendingTasks()
	if err == nil {
		t.Error("Expected error for database failure")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDatabaseService_GetTasksByExpressionID_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	service := &DatabaseService{db: db}

	rows := sqlmock.NewRows([]string{"id", "expression_id", "arg1", "arg2", "operation", "operation_time", "status", "result", "created_at", "updated_at"}).
		AddRow("task-id-1", "expr-id", "2", "2", "+", 1000, "pending", nil, time.Now(), time.Now()).
		AddRow("task-id-2", "expr-id", "3", "3", "+", 1000, "completed", &[]float64{6.0}[0], time.Now(), time.Now())

	mock.ExpectQuery("SELECT id, expression_id, arg1, arg2, operation, operation_time, status, result, created_at, updated_at FROM tasks WHERE expression_id = \\? ORDER BY created_at ASC").
		WithArgs("expr-id").
		WillReturnRows(rows)

	tasks, err := service.GetTasksByExpressionID("expr-id")
	if err != nil {
		t.Fatalf("Failed to get tasks by expression ID: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	mock.ExpectQuery("SELECT id, expression_id, arg1, arg2, operation, operation_time, status, result, created_at, updated_at FROM tasks WHERE expression_id = \\? ORDER BY created_at ASC").
		WithArgs("expr-id").
		WillReturnError(errors.New("database error"))

	_, err = service.GetTasksByExpressionID("expr-id")
	if err == nil {
		t.Error("Expected error for database failure")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDatabaseService_Close_Mock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	service := &DatabaseService{db: db}

	mock.ExpectClose()

	err = service.Close()
	if err != nil {
		t.Errorf("Failed to close database: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
