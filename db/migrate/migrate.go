// Inspired by https://github.com/fragmenta/fragmenta/blob/master/migrate.go
package main

import (
	"fmt"
	"log"
	"os"
	"slices"

	"path/filepath"
	"sort"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

var (
	database = os.Getenv("DB")
)

var conn *sqlite.Conn

const (
	migrationsPath = "./db/migrations"
)

func Conn() *sqlite.Conn {
	if conn != nil {
		return conn
	}
	dbPath := database
	if dbPath == "" {
		dbPath = "app.db"
	}
	conn, err := sqlite.OpenConn(dbPath, sqlite.OpenReadWrite|sqlite.OpenCreate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn
}

func readMetadata() []string {
	var migrations []string

	sql := "select name from _migrations order by name asc;"

	stmt, _, err := Conn().PrepareTransient(sql)
	if err != nil {
		log.Printf("Database ERROR preparing statement: %v", err)
		return migrations
	}
	defer stmt.Finalize()

	for {
		row, err := stmt.Step()
		if err != nil {
			log.Printf("Database ERROR %s", err)
			return migrations
		}
		if !row {
			break
		}
		migrations = append(migrations, stmt.ColumnText(0))
	}

	return migrations
}

func Migrate(migrationName string) {

	var migrations []string
	var completed []string

	// Get a list of migration files
	files, err := filepath.Glob(migrationsPath + "/*.sql")
	if err != nil {
		log.Printf("Error running restore %s", err)
		return
	}

	// Sort the list alphabetically
	sort.Strings(files)

	conn = Conn()
	defer conn.Close()

	// Create migrations table if it doesn't exist
	createTableSQL := `CREATE TABLE IF NOT EXISTS _migrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	err = sqlitex.ExecuteTransient(conn, createTableSQL, nil)
	if err != nil {
		log.Printf("Error creating migrations table: %v", err)
		return
	}

	migrations = readMetadata()

	for _, file := range files {
		filename := filepath.Base(file)

		if migrationName != "" && filename != migrationName {
			continue
		}

		if !contains(filename, migrations) {
			log.Printf("Running migration %s", filename)

			// Read the SQL file content
			sqlContent, err := os.ReadFile(file)
			if err != nil {
				log.Printf("ERROR reading migration file %s: %v", filename, err)
				log.Printf("All further migrations cancelled\n\n")
				break
			}

			// Execute the SQL content against the SQLite database
			err = sqlitex.ExecuteScript(conn, string(sqlContent), nil)
			if err != nil {
				log.Printf("ERROR executing migration %s: %v", filename, err)
				log.Printf("All further migrations cancelled\n\n")
				break
			}

			completed = append(completed, filename)
			log.Printf("Completed migration %s\n\n", filename)
		}
	}

	if len(completed) > 0 {
		writeMetadata(completed)
		log.Printf("Migrations complete up to migration %s on db %s\n\n", completed[len(completed)-1], database)
	} else {
		log.Printf("No migrations to perform at path %s\n\n", migrationsPath)
	}
}

// writeMetadata writes a new row in the fragmenta_metadata table to record our action
func writeMetadata(migrations []string) {

	for _, m := range migrations {
		sql := "Insert into _migrations(name) VALUES(?);"
		err := sqlitex.ExecuteTransient(Conn(), sql, &sqlitex.ExecOptions{
			Args: []any{m},
		})
		if err != nil {
			log.Printf("Database ERROR %s", err)
		}
	}

}

// contains checks whether an array of strings contains a string
func contains(s string, a []string) bool {
	return slices.Contains(a, s)
}
