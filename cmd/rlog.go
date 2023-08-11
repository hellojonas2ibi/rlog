package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hellojonas/rlog/internal/rlog"
)

func main() {
	logDir := os.Getenv("RLOG_DIR")
	intervalStr := os.Getenv("RLOG_INTERVAL")
	bufferSizeStr := os.Getenv("RLOG_BUFFER_SIZE")

	if logDir == "" {
		logDir = os.Getenv("$HOME/.rlog/logs")
	}

	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		interval = 10
	}

	bufferSize, err := strconv.Atoi(bufferSizeStr)
	if err != nil {
		bufferSize = 20
	}

	logger := rlog.New(logDir, bufferSize, time.Duration(interval)*time.Second)
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Mount("/api", rlog.Routes(logger))

	http.ListenAndServe(":8080", r)
}
