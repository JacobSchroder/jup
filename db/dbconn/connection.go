package dbconn

import (
	"context"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func NewConnectionPool() (*sqlitex.Pool, error) {
	return sqlitex.NewPool("local-data/app.db", sqlitex.PoolOptions{
		PoolSize: 10,
		Flags:    sqlite.OpenReadWrite | sqlite.OpenCreate | sqlite.OpenWAL,
	})
}

func GetConnection(pool *sqlitex.Pool, ctx context.Context) (*sqlite.Conn, error) {
	return pool.Take(ctx)
}

func PutConnection(pool *sqlitex.Pool, conn *sqlite.Conn) {
	pool.Put(conn)
}
