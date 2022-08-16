package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func rootHandler(api http.Handler) chi.Router {
	handler := chi.NewRouter()

	handler.Mount("/api", api)

	return handler
}

func apiHandler(varlog http.Handler) chi.Router {
	handler := chi.NewRouter()

	handler.Mount("/varlog", varlog)

	return handler
}
