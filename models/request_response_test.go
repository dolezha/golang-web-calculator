package models

import (
	"encoding/json"
	"testing"
)

func TestRequestBody_JSONSerialization(t *testing.T) {
	req := RequestBody{
		Expression: "2+3*4",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	jsonStr := string(data)
	if !contains(jsonStr, "2+3*4") {
		t.Error("Expression should be serialized to JSON")
	}

	var unmarshaled RequestBody
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal request: %v", err)
	}

	if unmarshaled.Expression != req.Expression {
		t.Errorf("Expected expression %s, got %s", req.Expression, unmarshaled.Expression)
	}
}

func TestResponseBody_JSONSerialization_WithResult(t *testing.T) {
	result := 42.0
	resp := ResponseBody{
		Result: &result,
		Error:  nil,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	jsonStr := string(data)
	if !contains(jsonStr, "42") {
		t.Error("Result should be serialized to JSON")
	}
	if contains(jsonStr, "error") {
		t.Error("Error should be omitted when nil")
	}

	var unmarshaled ResponseBody
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if unmarshaled.Result == nil || *unmarshaled.Result != *resp.Result {
		t.Errorf("Expected result %v, got %v", resp.Result, unmarshaled.Result)
	}
}

func TestResponseBody_JSONSerialization_WithError(t *testing.T) {
	errorMsg := "invalid expression"
	resp := ResponseBody{
		Result: nil,
		Error:  &errorMsg,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	jsonStr := string(data)
	if !contains(jsonStr, "invalid expression") {
		t.Error("Error should be serialized to JSON")
	}
	if contains(jsonStr, "result") {
		t.Error("Result should be omitted when nil")
	}

	var unmarshaled ResponseBody
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if unmarshaled.Error == nil || *unmarshaled.Error != *resp.Error {
		t.Errorf("Expected error %v, got %v", resp.Error, unmarshaled.Error)
	}
}

func TestResponseBody_JSONSerialization_Empty(t *testing.T) {
	resp := ResponseBody{
		Result: nil,
		Error:  nil,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	jsonStr := string(data)
	if contains(jsonStr, "result") {
		t.Error("Result should be omitted when nil")
	}
	if contains(jsonStr, "error") {
		t.Error("Error should be omitted when nil")
	}

	var unmarshaled ResponseBody
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if unmarshaled.Result != nil {
		t.Errorf("Expected nil result, got %v", unmarshaled.Result)
	}
	if unmarshaled.Error != nil {
		t.Errorf("Expected nil error, got %v", unmarshaled.Error)
	}
}
