package logger

import (
	"os"
	"testing"
)

func TestInitFromEnv(t *testing.T) {
	os.Setenv("LOG_FORMAT", "text")
	os.Setenv("LOG_LEVEL", "debug")
	defer os.Unsetenv("LOG_FORMAT")
	defer os.Unsetenv("LOG_LEVEL")

	if err := InitFromEnv(); err != nil {
		t.Fatalf("InitFromEnv failed: %v", err)
	}
}
