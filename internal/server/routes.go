package server

import (
	"net/http"

	"github.com/JacobSchroder/jup/internal/handlers"
)

func AddRoutes(
	mux *http.ServeMux,
) {

	mux.HandleFunc("GET /{$}", handlers.HandleGetIndex)

	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/",   http.StripPrefix("/static/", fileServer))
}
