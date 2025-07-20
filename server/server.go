package server

import (
	"net/http"

	"github.com/JacobSchroder/jup/server/routes"
)

type Config struct {
	port int
}

func Server() http.Handler {
	mux := http.NewServeMux()

	routes.AddRoutes(mux)
	var handler http.Handler = mux

	return handler
}
