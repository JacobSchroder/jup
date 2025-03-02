package handlers

import (
	"net/http"

	"github.com/JacobSchroder/jup/internal/templates"
)

func HandleGetPostCommentForm(w http.ResponseWriter, r *http.Request){

	err := templates.CommentForm().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}