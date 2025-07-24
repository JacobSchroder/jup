package server

import (
	"net/http"

	"github.com/JacobSchroder/jup/db/dbconn"
	"github.com/JacobSchroder/jup/server/di"
	"github.com/JacobSchroder/jup/server/routes"
)

type Config struct {
	port int
}

func Server() http.Handler {
	mux := http.NewServeMux()

	pool, err := dbconn.NewConnectionPool()
	if err != nil {
		panic(err)
	}

	app := &di.App{DB: pool}

	routes.AddRoutes(mux, app)
	var handler http.Handler = mux

	return handler
}
