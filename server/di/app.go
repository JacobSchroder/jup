package di

import (
	"zombiezen.com/go/sqlite/sqlitex"
)

type App struct {
	DB *sqlitex.Pool
	WebSocketHandler interface{}
}
