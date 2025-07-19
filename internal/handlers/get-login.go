package handlers

import (
	"net/http"

	"github.com/JacobSchroder/jup/internal/templates"
)

func HandleGetLogin(w http.ResponseWriter, r *http.Request) {
	errorMsg := r.URL.Query().Get("error")
	component := templates.Layout(templates.Login(errorMsg), "Login")
	component.Render(r.Context(), w)
}
