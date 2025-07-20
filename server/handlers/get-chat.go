package handlers

import (
	"net/http"

	"github.com/JacobSchroder/jup/pages"
)

func HandleGetChat(w http.ResponseWriter, r *http.Request) {
	content := pages.Chat()
	err := content.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
