package handlers

import (
	"net/http"

	"github.com/JacobSchroder/jup/pages"
)

func HandleGetIndex(w http.ResponseWriter, r *http.Request) {
	content := pages.Index()
	err := content.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
