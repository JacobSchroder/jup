package server

import (
	"net/http"
)

type Config struct {
	port int
}



func Server() http.Handler {
	mux := http.NewServeMux()

	AddRoutes(mux)
	var handler http.Handler = mux

	return handler
}