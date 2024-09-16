package main

import (
	"github.com/sebastianreh/user-balance-api/cmd/httpserver"
	_ "github.com/sebastianreh/user-balance-api/docs/swagger"
	"github.com/sebastianreh/user-balance-api/internal/container"
	"github.com/sebastianreh/user-balance-api/internal/interfaces/middlewares"
)

//	@title			User Balance API
//	@version		0.0.1
//	@description	This is the API documentation for the User Balance service.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8000
// @BasePath	/user-balance-api/
func main() {
	dependencies := container.Build()
	server := httpserver.NewServer(dependencies)
	middlewares.AddMiddlewares(server.Server, middlewares.WithRecover())
	server.Routes()
	server.SetErrorHandler(middlewares.HTTPErrorHandler)
	server.Start()
}
