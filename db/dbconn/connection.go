package dbconn

import (
	"context"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func NewConnectionPool() (*sqlitex.Pool, error) {
	return sqlitex.NewPool("local-data/app.db", sqlitex.PoolOptions{})
}

func GetConnection(pool *sqlitex.Pool) (*sqlite.Conn, error) {
	return pool.Take(context.TODO())
}

func PutConnection(pool *sqlitex.Pool, conn *sqlite.Conn) {
	pool.Put(conn)
}
