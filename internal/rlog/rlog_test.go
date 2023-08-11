package rlog

import (
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestOpenLogFile(t *testing.T) {
	logDir := filepath.Join("logs")
	logger := New(logDir, 0, 100)
	file, err := logger.logFile("logger1")

	if err != nil {
		t.Fatalf("error opening log file. %v", err)
	}

	logger.Close()
	file.Close()
}

func TestLogMultipleFiles(t *testing.T) {
	logDir := filepath.Join("logs")
	logger := New(logDir, 100, time.Duration(5)*time.Second)
	n := 15
	lines := 10000
	tenants := make([]string, n)

	for i := 0; i < n; i++ {
		tenants[i] = fmt.Sprintf("tenant%d", i+1)
	}

	wg := sync.WaitGroup{}

	for _, tenant := range tenants {
		wg.Add(1)
		go func(tenant string) {
			for i := 0; i < lines; i++ {
				entry := Entry{
					identifier: tenant,
					Message:    fmt.Sprintf("[%s][%-4d] %s", tenant, i, "This is a log entry."),
				}
				logger.Submit(entry)
				time.Sleep(time.Duration(10) * time.Millisecond)
			}
			wg.Done()
		}(tenant)
	}

	wg.Wait()

	for {
		select {
		case <-logger.Exit:
			break
		}
		fmt.Println("---------> [WAITING]")
	}
}
