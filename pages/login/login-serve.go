package page_login

import (
	"net/http"
)

func HandleGetLogin(w http.ResponseWriter, r *http.Request) {
	errorMsg := r.URL.Query().Get("error")
	component := Login(errorMsg)
	component.Render(r.Context(), w)
}
