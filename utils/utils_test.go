package utils

import (
	"net/http/httptest"
	"testing"
)

func TestToPtr(t *testing.T) {
	value := 42
	ptr := ToPtr(value)

	if ptr == nil {
		t.Error("Expected non-nil pointer")
	}

	if *ptr != value {
		t.Errorf("Expected %d, got %d", value, *ptr)
	}
}

func TestRespondWithJSON(t *testing.T) {
	rr := httptest.NewRecorder()

	data := map[string]string{"message": "test"}
	err := RespondWithJSON(rr, data, 200)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if rr.Code != 200 {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	expected := `{"message":"test"}` + "\n"
	if rr.Body.String() != expected {
		t.Errorf("Expected %s, got %s", expected, rr.Body.String())
	}
}
