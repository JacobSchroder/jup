package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/JacobSchroder/jup/server/di"
	"github.com/JacobSchroder/jup/server/utils"
	"github.com/JacobSchroder/jup/templates/issue"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func HandlePostIssue(app *di.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		title := r.PostFormValue("title")
		description := r.PostFormValue("description")

		if title == "" || description == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		id, err := utils.GenerateID(utils.PrefixIssue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		conn, err := app.DB.Take(context.TODO())
		if err != nil {
			http.Error(w, "Unable to connect to database", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		err = sqlitex.ExecuteTransient(conn, "INSERT INTO issues (id, title, description) VALUES (?, ?, ?);", &sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				fmt.Println(stmt.ColumnText(0))
				return nil
			},
			Args: []any{id, title, description},
		})

		if err != nil {
			http.Error(w, "Failed to create issue", http.StatusInternalServerError)
			return
		}

		err = issue.Issue(issue.IssueProps{Id: id, Title: title, Description: description}).Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
	}
}
