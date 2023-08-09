package rlog

import (
	"path/filepath"
	"testing"
)

func TestOpenLogFile(t *testing.T) {
	logDir := filepath.Join(".rlog")
	logger := New(logDir)
	file, err := logger.logFile("logger1")

	if err != nil {
		t.Fatalf("error opening log file. %v", err)
	}

	file.Close()
}
