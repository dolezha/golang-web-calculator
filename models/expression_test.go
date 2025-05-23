package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestExpressionStatus_Constants(t *testing.T) {
	if StatusPending != "pending" {
		t.Errorf("Expected StatusPending to be 'pending', got %s", StatusPending)
	}
	if StatusComputing != "computing" {
		t.Errorf("Expected StatusComputing to be 'computing', got %s", StatusComputing)
	}
	if StatusDone != "done" {
		t.Errorf("Expected StatusDone to be 'done', got %s", StatusDone)
	}
}

func TestExpression_JSONSerialization(t *testing.T) {
	result := 42.0
	expr := Expression{
		ID:         "test-id",
		UserID:     1,
		Expression: "2+2",
		Status:     StatusDone,
		Result:     &result,
		CreatedAt:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:  time.Date(2023, 1, 1, 0, 0, 1, 0, time.UTC),
	}

	data, err := json.Marshal(expr)
	if err != nil {
		t.Fatalf("Failed to marshal expression: %v", err)
	}

	jsonStr := string(data)
	if !contains(jsonStr, "test-id") {
		t.Error("ID should be serialized to JSON")
	}
	if !contains(jsonStr, "2+2") {
		t.Error("Expression should be serialized to JSON")
	}
	if !contains(jsonStr, "done") {
		t.Error("Status should be serialized to JSON")
	}

	var unmarshaled Expression
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal expression: %v", err)
	}

	if unmarshaled.ID != expr.ID {
		t.Errorf("Expected ID %s, got %s", expr.ID, unmarshaled.ID)
	}
	if unmarshaled.Status != expr.Status {
		t.Errorf("Expected status %s, got %s", expr.Status, unmarshaled.Status)
	}
	if unmarshaled.Result == nil || *unmarshaled.Result != *expr.Result {
		t.Errorf("Expected result %v, got %v", expr.Result, unmarshaled.Result)
	}
}

func TestExpression_JSONSerialization_NilResult(t *testing.T) {
	expr := Expression{
		ID:         "test-id",
		UserID:     1,
		Expression: "2+2",
		Status:     StatusPending,
		Result:     nil,
		CreatedAt:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:  time.Date(2023, 1, 1, 0, 0, 1, 0, time.UTC),
	}

	data, err := json.Marshal(expr)
	if err != nil {
		t.Fatalf("Failed to marshal expression: %v", err)
	}

	jsonStr := string(data)
	if contains(jsonStr, "result") {
		t.Error("Result should be omitted when nil")
	}

	var unmarshaled Expression
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal expression: %v", err)
	}

	if unmarshaled.Result != nil {
		t.Errorf("Expected nil result, got %v", unmarshaled.Result)
	}
}
