package routes

import (
	"net/http"

	page_issues_list "github.com/JacobSchroder/jup/pages/issues-list"
	page_login "github.com/JacobSchroder/jup/pages/login"
	page_websocket_test "github.com/JacobSchroder/jup/pages/websocket-test"

	"github.com/JacobSchroder/jup/server/di"
	"github.com/JacobSchroder/jup/server/handlers"
)

func AddRoutes(
	mux *http.ServeMux,
	app *di.App,
) {

	mux.HandleFunc("GET /{$}", page_issues_list.HandleGetIssues(app))

	mux.HandleFunc("GET /login", page_login.HandleGetLogin)
	mux.HandleFunc("POST /login", handlers.HandlePostLogin)

	// WebSocket test page
	mux.HandleFunc("GET /ws-test", page_websocket_test.HandleGetWebSocketTest)

	mux.HandleFunc("POST /comment", handlers.HandlePostComment)
	mux.HandleFunc("GET /comment-form", handlers.HandleGetPostCommentForm)

	mux.HandleFunc("POST /issues", handlers.HandlePostIssue(app))
	mux.HandleFunc("DELETE /issues/{issueId}", handlers.HandleDeleteIssue(app))

	// WebSocket endpoint
	if wsHandler, ok := app.WebSocketHandler.(*handlers.WebSocketHandler); ok {
		mux.HandleFunc("GET /ws", wsHandler.HandleWebSocket())
	}

	fileServer := http.FileServer(http.Dir("assets"))
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", fileServer))
}
