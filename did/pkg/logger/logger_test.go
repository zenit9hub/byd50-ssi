package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFuncName(t *testing.T) {
	name := helperFuncName()
	if !strings.Contains(name, "helperFuncName") {
		t.Fatalf("unexpected func name: %s", name)
	}
}

func helperFuncName() string {
	return FuncName()
}

func TestGetFunctionName(t *testing.T) {
	name := GetFunctionName(sampleFunction)
	if !strings.Contains(name, "sampleFunction") {
		t.Fatalf("unexpected function name: %s", name)
	}
}

func sampleFunction() {
}

func TestFuncStartEnd(t *testing.T) {
	FuncStart()
	FuncEnd()
}

func TestLoggerFileOutput(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test.log")

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Chdir(cwd)
	}()

	SetFileName(logPath)
	Log("hello")
	Debug("debug")
	Error("error")

	closeFile()

	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Fatal("expected log file to have content")
	}

	if getCallingFunctionName() == "" {
		t.Fatal("expected calling function name")
	}

	lastLogDate = "2000-01-01"
	fileName = "old.log"
	initialized = true
	Log("rotate")
	Debug("rotate-debug")
	Error("rotate-error")

	lastLogDate = "2000-01-01"
	fileName = "old-debug.log"
	initialized = true
	Debug("rotate-debug-2")

	lastLogDate = "2000-01-01"
	fileName = "old-error.log"
	initialized = true
	Error("rotate-error-2")
}
