package handlers

import (
	"net/http"

	"github.com/JacobSchroder/jup/internal/templates"
)

func HandleGetIndex(w http.ResponseWriter, r *http.Request){
	content := templates.Index();
	err := templates.Layout(content,  "Hello, world!").Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	return
}