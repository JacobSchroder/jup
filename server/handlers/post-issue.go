package handlers

import (
	"net/http"

	"github.com/JacobSchroder/jup/templates/issue"
)

func HandlePostIssue(w http.ResponseWriter, r *http.Request) {

	title := r.PostFormValue("title")
	description := r.PostFormValue("description")

	if title == "" || description == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	err := issue.Issue(issue.IssueProps{Title: title, Description: description}).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

}
