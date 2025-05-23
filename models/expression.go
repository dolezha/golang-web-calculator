package models

import "time"

type ExpressionStatus string

const (
	StatusPending   ExpressionStatus = "pending"
	StatusComputing ExpressionStatus = "computing"
	StatusDone      ExpressionStatus = "done"
)

type Expression struct {
	ID         string           `json:"id" db:"id"`
	UserID     int              `json:"user_id" db:"user_id"`
	Expression string           `json:"expression" db:"expression"`
	Status     ExpressionStatus `json:"status" db:"status"`
	Result     *float64         `json:"result,omitempty" db:"result"`
	CreatedAt  time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at" db:"updated_at"`
}
