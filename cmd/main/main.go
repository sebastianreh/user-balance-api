package main

import (
	"github.com/sebastianreh/user-balance-api/cmd/httpserver"
	"github.com/sebastianreh/user-balance-api/internal/container"
	"github.com/sebastianreh/user-balance-api/internal/interfaces/middlewares"
)

func main() {
	dependencies := container.Build()
	server := httpserver.NewServer(dependencies)
	middlewares.AddMiddlewares(server.Server, middlewares.WithLogger(dependencies.Config),
		middlewares.WithRecover(), middlewares.WithRecover())
	server.Routes()
	server.SetErrorHandler(middlewares.HTTPErrorHandler)
	server.Start()
}
