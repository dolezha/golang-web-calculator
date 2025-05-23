package main

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	result := getEnv("TEST_VAR", "default")
	if result != "test_value" {
		t.Errorf("Expected test_value, got %s", result)
	}

	result = getEnv("NON_EXISTING_VAR", "default")
	if result != "default" {
		t.Errorf("Expected default, got %s", result)
	}

	os.Setenv("EMPTY_VAR", "")
	defer os.Unsetenv("EMPTY_VAR")

	result = getEnv("EMPTY_VAR", "default")
	if result != "default" {
		t.Errorf("Expected default for empty env var, got %s", result)
	}
}

func TestGetEnvVariousScenarios(t *testing.T) {
	testCases := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		setEnv       bool
		expected     string
	}{
		{"existing_env_var", "TEST_KEY_1", "default1", "value1", true, "value1"},
		{"non_existing_env_var", "TEST_KEY_2", "default2", "", false, "default2"},
		{"empty_env_var", "TEST_KEY_3", "default3", "", true, "default3"},
		{"whitespace_env_var", "TEST_KEY_4", "default4", "   ", true, "   "},
		{"special_chars", "TEST_KEY_5", "default5", "value@#$%", true, "value@#$%"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				os.Setenv(tc.key, tc.envValue)
				defer os.Unsetenv(tc.key)
			}

			result := getEnv(tc.key, tc.defaultValue)
			if result != tc.expected {
				t.Errorf("getEnv(%s, %s) = %s, expected %s", tc.key, tc.defaultValue, result, tc.expected)
			}
		})
	}
}

func TestMain(t *testing.T) {
	if os.Getenv("RUN_MAIN_TEST") == "1" {
		t.Skip("Skipping main execution due to CGO dependency")
	}

}

func TestMainPackageExists(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping main package test in short mode")
	}

	longValue := string(make([]byte, 1000))
	for i := range longValue {
		longValue = longValue[:i] + "a" + longValue[i+1:]
	}

	os.Setenv("LONG_VAR", longValue)
	defer os.Unsetenv("LONG_VAR")

	result := getEnv("LONG_VAR", "default")
	if result != longValue {
		t.Error("getEnv should handle long values correctly")
	}

	unicodeValue := "测试值🚀"
	os.Setenv("UNICODE_VAR", unicodeValue)
	defer os.Unsetenv("UNICODE_VAR")

	result = getEnv("UNICODE_VAR", "default")
	if result != unicodeValue {
		t.Error("getEnv should handle unicode values correctly")
	}
}
