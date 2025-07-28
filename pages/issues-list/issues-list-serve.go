package page_issues_list

import (
	"net/http"

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
				issues = append(issues, issue.IssueProps{Id: id, Title: title, Description: description})
				return nil
			},
		})

		err = IssuesList(IssuesListProps{Issues: issues}).Render(r.Context(), w)

	}
}
