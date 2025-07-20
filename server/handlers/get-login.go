package handlers

import (
	"net/http"

	"github.com/JacobSchroder/jup/pages"
)

func HandleGetLogin(w http.ResponseWriter, r *http.Request) {
	errorMsg := r.URL.Query().Get("error")
	component := pages.Login(errorMsg)
	component.Render(r.Context(), w)
}
