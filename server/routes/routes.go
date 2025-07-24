package routes

import (
	"net/http"

	"github.com/JacobSchroder/jup/server/di"
	"github.com/JacobSchroder/jup/server/handlers"
)

func AddRoutes(
	mux *http.ServeMux,
	app *di.App,
) {

	mux.HandleFunc("GET /{$}", handlers.HandleGetIssues(app))

	mux.HandleFunc("GET /login", handlers.HandleGetLogin)
	mux.HandleFunc("POST /login", handlers.HandlePostLogin)

	mux.HandleFunc("POST /comment", handlers.HandlePostComment)
	mux.HandleFunc("GET /comment-form", handlers.HandleGetPostCommentForm)

	mux.HandleFunc("POST /issues", handlers.HandlePostIssue(app))

	fileServer := http.FileServer(http.Dir("assets"))
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", fileServer))
}
