package models

import "time"

type Task struct {
	ID            string    `json:"id" db:"id"`
	ExpressionID  string    `json:"expression_id" db:"expression_id"`
	Arg1          string    `json:"arg1" db:"arg1"`
	Arg2          string    `json:"arg2" db:"arg2"`
	Operation     string    `json:"operation" db:"operation"`
	OperationTime int64     `json:"operation_time" db:"operation_time"`
	Status        string    `json:"status" db:"status"`
	Result        *float64  `json:"result,omitempty" db:"result"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
