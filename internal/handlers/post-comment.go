package handlers

import (
	"log/slog"
	"net/http"

	templates "github.com/JacobSchroder/jup/internal/templates/form"
)

func HandlePostComment(w http.ResponseWriter, r *http.Request){

	comment := r.PostFormValue("comment")

	slog.Info("Formvalues", "comment", comment)

	err := templates.CommentPosted(comment).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

}