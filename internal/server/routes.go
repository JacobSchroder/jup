package server

import (
	"net/http"

	"github.com/JacobSchroder/jup/internal/handlers"
)

func AddRoutes(
	mux *http.ServeMux,
) {

	mux.HandleFunc("GET /{$}", handlers.HandleGetIndex)

	mux.HandleFunc("POST /comment", handlers.HandlePostComment)
	mux.HandleFunc("GET /comment-form", handlers.HandleGetPostCommentForm)

	fileServer := http.FileServer(http.Dir("assets"))
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", fileServer))
}
