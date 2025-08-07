package websocket_test

import (
	"net/http"
)

func HandleGetWebSocketTest(w http.ResponseWriter, r *http.Request) {
	component := WebSocketTest()
	component.Render(r.Context(), w)
}
