package models

type ExpressionStatus string

const (
	StatusPending   ExpressionStatus = "pending"
	StatusComputing ExpressionStatus = "computing"
	StatusDone      ExpressionStatus = "done"
)

type Expression struct {
	ID         string           `json:"id"`
	Expression string           `json:"expression"`
	Status     ExpressionStatus `json:"status"`
	Result     *float64         `json:"result,omitempty"`
}
