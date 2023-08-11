package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/hellojonas/rlog/internal/rlog"
)

func TestPostingLogs(t *testing.T) {
	baseUrl := "http://localhost:8080/api/logs"
	n := 10
	lines := 10000
	tenants := make([]string, n)

	for i := 0; i < n; i++ {
		tenants[i] = fmt.Sprintf("tenant%d", i+1)
	}

	wg := sync.WaitGroup{}

	s := rand.NewSource(time.Now().Unix())
	_ = rand.New(s)

	for _, tenant := range tenants {
		wg.Add(1)
		go func(tenant string) {
			url := baseUrl + "/" + tenant
			for i := 0; i < lines; i++ {
				entry := rlog.Entry{
					Message: fmt.Sprintf("[%s][%05d] %s", tenant, i, "This is a log entry."),
				}
				buf := bytes.NewBuffer(make([]byte, 0))
				encoder := json.NewEncoder(buf)

				if err := encoder.Encode(&entry); err != nil {
					fmt.Printf("[EROR] error encoding log entry %v\n", err)
				}

				time.Sleep(time.Duration(50) * time.Millisecond)
				http.Post(url, "application/json", buf)
			}
			wg.Done()
		}(tenant)
	}

	wg.Wait()
}
