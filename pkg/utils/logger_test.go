package utils

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// captureStdout captures everything written to os.Stdout during the execution of f.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	os.Stdout = w

	outChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outChan <- buf.String()
	}()

	f()

	_ = w.Close()
	os.Stdout = oldStdout
	output := <-outChan
	_ = r.Close()

	return output
}

func TestLogger_WithColor(t *testing.T) {
	orig := colorEnabled
	defer func() { colorEnabled = orig }()
	colorEnabled = true

	out := captureStdout(t, func() {
		Error("database connection failed")
		Success("build succeeded")
		Warning("deprecated flag used")
		Info("listening on port 3030")
		Log("plain log message")
	})

	if !strings.Contains(out, ColorRed+"ERROR: database connection failed"+ColorReset) {
		t.Errorf("expected colored Error output, got: %q", out)
	}
	if !strings.Contains(out, ColorGreen+"build succeeded"+ColorReset) {
		t.Errorf("expected colored Success output, got: %q", out)
	}
	if !strings.Contains(out, ColorYellow+"WARNING: deprecated flag used"+ColorReset) {
		t.Errorf("expected colored Warning output, got: %q", out)
	}
	if !strings.Contains(out, ColorBlue+"INFO: listening on port 3030"+ColorReset) {
		t.Errorf("expected colored Info output, got: %q", out)
	}
	if !strings.Contains(out, "plain log message\n") {
		t.Errorf("expected plain Log output, got: %q", out)
	}
}

func TestLogger_WithoutColor(t *testing.T) {
	orig := colorEnabled
	defer func() { colorEnabled = orig }()
	colorEnabled = false

	out := captureStdout(t, func() {
		Error("database connection failed")
		Success("build succeeded")
		Warning("deprecated flag used")
		Info("listening on port 3030")
		Log("plain log message")
	})

	// Ensure no ANSI escape sequences exist in the output
	if strings.Contains(out, "\033[") {
		t.Errorf("expected no ANSI escape sequences, but found some in: %q", out)
	}

	if !strings.Contains(out, "ERROR: database connection failed\n") {
		t.Errorf("expected plain Error output, got: %q", out)
	}
	if !strings.Contains(out, "build succeeded\n") {
		t.Errorf("expected plain Success output, got: %q", out)
	}
	if !strings.Contains(out, "WARNING: deprecated flag used\n") {
		t.Errorf("expected plain Warning output, got: %q", out)
	}
	if !strings.Contains(out, "INFO: listening on port 3030\n") {
		t.Errorf("expected plain Info output, got: %q", out)
	}
	if !strings.Contains(out, "plain log message\n") {
		t.Errorf("expected plain Log output, got: %q", out)
	}
}

func TestLoggerStruct_WithColor(t *testing.T) {
	orig := colorEnabled
	defer func() { colorEnabled = orig }()
	colorEnabled = true

	logger := NewLogger("TEST")
	out := captureStdout(t, func() {
		logger.Error("struct error")
		logger.Success("struct success")
		logger.Warning("struct warning")
		logger.Info("struct info")
		logger.Debug("struct debug")
	})

	if !strings.Contains(out, ColorRed) || !strings.Contains(out, "TEST ERROR: struct error"+ColorReset) {
		t.Errorf("expected colored struct Error output, got: %q", out)
	}
	if !strings.Contains(out, ColorGreen) || !strings.Contains(out, "TEST SUCCESS: struct success"+ColorReset) {
		t.Errorf("expected colored struct Success output, got: %q", out)
	}
	if !strings.Contains(out, ColorYellow) || !strings.Contains(out, "TEST WARNING: struct warning"+ColorReset) {
		t.Errorf("expected colored struct Warning output, got: %q", out)
	}
	if !strings.Contains(out, ColorBlue) || !strings.Contains(out, "TEST INFO: struct info"+ColorReset) {
		t.Errorf("expected colored struct Info output, got: %q", out)
	}
	if !strings.Contains(out, "TEST DEBUG: struct debug"+ColorReset) {
		t.Errorf("expected struct Debug output, got: %q", out)
	}
}

func TestLoggerStruct_WithoutColor(t *testing.T) {
	orig := colorEnabled
	defer func() { colorEnabled = orig }()
	colorEnabled = false

	logger := NewLogger("TEST")
	out := captureStdout(t, func() {
		logger.Error("struct error")
		logger.Success("struct success")
		logger.Warning("struct warning")
		logger.Info("struct info")
		logger.Debug("struct debug")
	})

	// Ensure no ANSI escape sequences exist in the output
	if strings.Contains(out, "\033[") {
		t.Errorf("expected no ANSI escape sequences in struct logger, but found some in: %q", out)
	}

	if !strings.Contains(out, "TEST ERROR: struct error\n") {
		t.Errorf("expected plain struct Error output, got: %q", out)
	}
	if !strings.Contains(out, "TEST SUCCESS: struct success\n") {
		t.Errorf("expected plain struct Success output, got: %q", out)
	}
	if !strings.Contains(out, "TEST WARNING: struct warning\n") {
		t.Errorf("expected plain struct Warning output, got: %q", out)
	}
	if !strings.Contains(out, "TEST INFO: struct info\n") {
		t.Errorf("expected plain struct Info output, got: %q", out)
	}
	if !strings.Contains(out, "TEST DEBUG: struct debug\n") {
		t.Errorf("expected plain struct Debug output, got: %q", out)
	}
}

func TestIsColorSupported_EnvironmentVariables(t *testing.T) {
	origNoColor := os.Getenv("NO_COLOR")
	origTerm := os.Getenv("TERM")
	defer func() {
		_ = os.Setenv("NO_COLOR", origNoColor)
		_ = os.Setenv("TERM", origTerm)
	}()

	// When NO_COLOR is set, isColorSupported must return false
	_ = os.Setenv("NO_COLOR", "1")
	if isColorSupported() {
		t.Errorf("expected isColorSupported() to be false when NO_COLOR is set")
	}

	_ = os.Unsetenv("NO_COLOR")
	_ = os.Setenv("TERM", "dumb")
	if isColorSupported() {
		t.Errorf("expected isColorSupported() to be false when TERM=dumb")
	}
}
