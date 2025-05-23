package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTask_JSONSerialization(t *testing.T) {
	result := 42.0
	task := Task{
		ID:            "task-1",
		ExpressionID:  "expr-1",
		Arg1:          "2",
		Arg2:          "3",
		Operation:     "+",
		OperationTime: 1000,
		Status:        "done",
		Result:        &result,
		CreatedAt:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:     time.Date(2023, 1, 1, 0, 0, 1, 0, time.UTC),
	}

	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	jsonStr := string(data)
	if !contains(jsonStr, "task-1") {
		t.Error("ID should be serialized to JSON")
	}
	if !contains(jsonStr, "expr-1") {
		t.Error("ExpressionID should be serialized to JSON")
	}
	if !contains(jsonStr, "+") {
		t.Error("Operation should be serialized to JSON")
	}

	var unmarshaled Task
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal task: %v", err)
	}

	if unmarshaled.ID != task.ID {
		t.Errorf("Expected ID %s, got %s", task.ID, unmarshaled.ID)
	}
	if unmarshaled.Operation != task.Operation {
		t.Errorf("Expected operation %s, got %s", task.Operation, unmarshaled.Operation)
	}
	if unmarshaled.Result == nil || *unmarshaled.Result != *task.Result {
		t.Errorf("Expected result %v, got %v", task.Result, unmarshaled.Result)
	}
}

func TestTask_JSONSerialization_NilResult(t *testing.T) {
	task := Task{
		ID:            "task-1",
		ExpressionID:  "expr-1",
		Arg1:          "2",
		Arg2:          "3",
		Operation:     "+",
		OperationTime: 1000,
		Status:        "pending",
		Result:        nil,
		CreatedAt:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:     time.Date(2023, 1, 1, 0, 0, 1, 0, time.UTC),
	}

	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	jsonStr := string(data)
	if contains(jsonStr, "result") {
		t.Error("Result should be omitted when nil")
	}

	var unmarshaled Task
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal task: %v", err)
	}

	if unmarshaled.Result != nil {
		t.Errorf("Expected nil result, got %v", unmarshaled.Result)
	}
}
