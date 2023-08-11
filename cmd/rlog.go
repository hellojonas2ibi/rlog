package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hellojonas/rlog/internal/rlog"
)

func main() {
	logDir := "logs"
	logger := rlog.New(logDir, 10, time.Duration(10)*time.Second)
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Mount("/api", rlog.Routes(logger))

	http.ListenAndServe(":8080", r)
}
