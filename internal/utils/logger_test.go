package utils

import (
	"bytes"
	"os"
	"testing"
)

func TestLoggerLevels(t *testing.T) {
	// Create a temporary file to capture output.
	f, err := os.CreateTemp(t.TempDir(), "logtest")
	if err != nil {
		t.Fatalf("CreateTemp error: %v", err)
	}
	defer f.Close()

	SetOutput(f)

	Info("info message")
	Debug("debug message")
	Error("error message")

	// Flush and read.
	f.Sync()
	data, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	out := string(data)
	if !contains(out, "INFO: info message") {
		t.Fatalf("expected INFO prefix, got %s", out)
	}
	if !contains(out, "DEBUG: debug message") {
		t.Fatalf("expected DEBUG prefix, got %s", out)
	}
	if !contains(out, "ERROR: error message") {
		t.Fatalf("expected ERROR prefix, got %s", out)
	}
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
