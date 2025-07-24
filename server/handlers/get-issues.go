package handlers

import (
	"fmt"
	"net/http"

	"github.com/JacobSchroder/jup/pages"
	"github.com/JacobSchroder/jup/server/di"
	"github.com/JacobSchroder/jup/templates/issue"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func HandleGetIssues(app *di.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		conn, err := app.DB.Take(r.Context())
		if err != nil {
			http.Error(w, "Unable to connect to database", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		var issues []issue.IssueProps
		sqlitex.ExecuteTransient(conn, "SELECT id, title, description FROM issues;", &sqlitex.ExecOptions{
			ResultFunc: func(stmt *sqlite.Stmt) error {
				id := stmt.ColumnText(0)
				title := stmt.ColumnText(1)
				description := stmt.ColumnText(2)
				fmt.Printf("ID: %s, Title: %s, Description: %s\n", id, title, description)
				issues = append(issues, issue.IssueProps{Id: id, Title: title, Description: description})
				return nil
			},
		})

		err = pages.Index(pages.IndexProps{Issues: issues}).Render(r.Context(), w)

	}
}
