package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/JacobSchroder/jup/server"
)

const (
	host = "localhost"
	port = "8082"
)

func run(ctx context.Context, w io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(w, nil))
	slog.SetDefault(logger)

	serv := server.Server()
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: serv,
	}

	err := httpServer.ListenAndServe()

	if err != nil {
		slog.Error(err.Error())
	}

	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
