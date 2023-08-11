package rlog

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes(logger *rlogger) chi.Router {
	r := chi.NewRouter()

	r.Post("/logs/{id}", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			id := chi.URLParam(r, "id")
			decoder := json.NewDecoder(r.Body)
			defer r.Body.Close()

			var entry Entry

			if err := decoder.Decode(&entry); err != nil {
				w.Write(nil)
				return
			}
			entry.identifier = id

			// logger.Submit(entry)
			w.Write(nil)
		}()
	})

	return r
}
