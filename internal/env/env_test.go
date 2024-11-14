package env_test

import (
	"os"
	"testing"

	"github.com/Crandel/gmail/internal/env"
)

func TestGetEnv(t *testing.T) {
	t.Parallel()
	// Set up a mock environment variable
	os.Setenv("TEST_ENV", "mock_value")

	// Test with an existing environment variable
	result := env.GetEnv("TEST_ENV", "default_value")
	if result != "mock_value" {
		t.Errorf("Expected 'mock_value', got '%s'", result)
	}

	// Test with a non-existing environment variable and default value
	result = env.GetEnv("NOT_EXISTING_ENV", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", result)
	}
}
