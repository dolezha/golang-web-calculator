package services

import (
	"calculator/models"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DatabaseService struct {
	db *sql.DB
}

func NewDatabaseService(dbPath string) (*DatabaseService, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	service := &DatabaseService{db: db}
	if err := service.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	return service, nil
}

func (ds *DatabaseService) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			login TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS expressions (
			id TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL,
			expression TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			result REAL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id)
		)`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			expression_id TEXT NOT NULL,
			arg1 TEXT NOT NULL,
			arg2 TEXT NOT NULL,
			operation TEXT NOT NULL,
			operation_time INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			result REAL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (expression_id) REFERENCES expressions (id)
		)`,
	}

	for _, query := range queries {
		if _, err := ds.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %v", err)
		}
	}

	return nil
}

func (ds *DatabaseService) CreateUser(login, passwordHash string) (*models.User, error) {
	query := `INSERT INTO users (login, password_hash) VALUES (?, ?)`
	result, err := ds.db.Exec(query, login, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user id: %v", err)
	}

	return &models.User{
		ID:           int(id),
		Login:        login,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}, nil
}

func (ds *DatabaseService) GetUserByLogin(login string) (*models.User, error) {
	query := `SELECT id, login, password_hash, created_at FROM users WHERE login = ?`
	row := ds.db.QueryRow(query, login)

	var user models.User
	err := row.Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return &user, nil
}

func (ds *DatabaseService) CreateExpression(expr *models.Expression) error {
	query := `INSERT INTO expressions (id, user_id, expression, status, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?)`
	_, err := ds.db.Exec(query, expr.ID, expr.UserID, expr.Expression, expr.Status,
		expr.CreatedAt, expr.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create expression: %v", err)
	}
	return nil
}

func (ds *DatabaseService) GetExpression(id string, userID int) (*models.Expression, error) {
	var query string
	var args []interface{}

	if userID == 0 {
		query = `SELECT id, user_id, expression, status, result, created_at, updated_at 
				 FROM expressions WHERE id = ?`
		args = []interface{}{id}
	} else {
		query = `SELECT id, user_id, expression, status, result, created_at, updated_at 
				 FROM expressions WHERE id = ? AND user_id = ?`
		args = []interface{}{id, userID}
	}

	row := ds.db.QueryRow(query, args...)

	var expr models.Expression
	err := row.Scan(&expr.ID, &expr.UserID, &expr.Expression, &expr.Status,
		&expr.Result, &expr.CreatedAt, &expr.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("expression not found")
		}
		return nil, fmt.Errorf("failed to get expression: %v", err)
	}

	return &expr, nil
}

func (ds *DatabaseService) UpdateExpression(expr *models.Expression) error {
	query := `UPDATE expressions SET status = ?, result = ?, updated_at = ? WHERE id = ?`
	_, err := ds.db.Exec(query, expr.Status, expr.Result, time.Now(), expr.ID)
	if err != nil {
		return fmt.Errorf("failed to update expression: %v", err)
	}
	return nil
}

func (ds *DatabaseService) GetUserExpressions(userID int) ([]*models.Expression, error) {
	query := `SELECT id, user_id, expression, status, result, created_at, updated_at 
			  FROM expressions WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := ds.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get expressions: %v", err)
	}
	defer rows.Close()

	var expressions []*models.Expression
	for rows.Next() {
		var expr models.Expression
		err := rows.Scan(&expr.ID, &expr.UserID, &expr.Expression, &expr.Status,
			&expr.Result, &expr.CreatedAt, &expr.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan expression: %v", err)
		}
		expressions = append(expressions, &expr)
	}

	return expressions, nil
}

func (ds *DatabaseService) CreateTask(task *models.Task) error {
	query := `INSERT INTO tasks (id, expression_id, arg1, arg2, operation, operation_time, status, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := ds.db.Exec(query, task.ID, task.ExpressionID, task.Arg1, task.Arg2,
		task.Operation, task.OperationTime, task.Status, task.CreatedAt, task.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create task: %v", err)
	}
	return nil
}

func (ds *DatabaseService) GetTask(id string) (*models.Task, error) {
	query := `SELECT id, expression_id, arg1, arg2, operation, operation_time, status, result, created_at, updated_at 
			  FROM tasks WHERE id = ?`
	row := ds.db.QueryRow(query, id)

	var task models.Task
	err := row.Scan(&task.ID, &task.ExpressionID, &task.Arg1, &task.Arg2,
		&task.Operation, &task.OperationTime, &task.Status, &task.Result,
		&task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to get task: %v", err)
	}

	return &task, nil
}

func (ds *DatabaseService) UpdateTask(task *models.Task) error {
	query := `UPDATE tasks SET status = ?, result = ?, updated_at = ? WHERE id = ?`
	_, err := ds.db.Exec(query, task.Status, task.Result, time.Now(), task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}
	return nil
}

func (ds *DatabaseService) GetPendingTasks() ([]*models.Task, error) {
	query := `SELECT id, expression_id, arg1, arg2, operation, operation_time, status, result, created_at, updated_at 
			  FROM tasks WHERE status = 'pending' ORDER BY created_at ASC`
	rows, err := ds.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending tasks: %v", err)
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.ExpressionID, &task.Arg1, &task.Arg2,
			&task.Operation, &task.OperationTime, &task.Status, &task.Result,
			&task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %v", err)
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (ds *DatabaseService) GetTasksByExpressionID(expressionID string) ([]*models.Task, error) {
	query := `SELECT id, expression_id, arg1, arg2, operation, operation_time, status, result, created_at, updated_at 
			  FROM tasks WHERE expression_id = ? ORDER BY created_at ASC`
	rows, err := ds.db.Query(query, expressionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %v", err)
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.ExpressionID, &task.Arg1, &task.Arg2,
			&task.Operation, &task.OperationTime, &task.Status, &task.Result,
			&task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %v", err)
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (ds *DatabaseService) Close() error {
	return ds.db.Close()
}
