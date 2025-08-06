package handlers

import (
	"net/http"

	"github.com/JacobSchroder/jup/server/di"
	"zombiezen.com/go/sqlite/sqlitex"
)

func HandleDeleteIssue(app *di.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.PathValue("issueId")

		conn, err := app.DB.Take(r.Context())
		if err != nil {
			http.Error(w, "Unable to connect to database", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		err = sqlitex.ExecuteTransient(conn, "DELETE FROM issues WHERE id = ?;", &sqlitex.ExecOptions{
			Args: []any{id},
		})

		if err != nil {
			http.Error(w, "Failed to delete issue", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(""))
	}
}
