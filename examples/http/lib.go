package main

import (
	"io"
	"net/http"
)

type ScopedService struct {
	counter int
}

type SingletonService struct {
	counter int
}

type TransientService struct {
	counter int
}

func getFromQuery(r *http.Request, key string, defaultValue string) string {
	if !r.URL.Query().Has(key) {
		return defaultValue
	}

	return r.URL.Query().Get(key)
}

func writeInternalError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, msg)
}
