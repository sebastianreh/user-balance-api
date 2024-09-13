package middlewares

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sebastianreh/user-balance-api/cmd/httpserver"
	"github.com/sebastianreh/user-balance-api/cmd/httpserver/exceptions"
	"github.com/sebastianreh/user-balance-api/internal/infrastructure/config"
	"net/http"
	"strings"

	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type Middleware func(echo *echo.Echo)

// AddMiddlewares build the middlewares of the server
func AddMiddlewares(server *echo.Echo, middlewares ...Middleware) {
	for _, middleware := range middlewares {
		middleware(server)
	}
}

func WithLogger(cfg config.Config) Middleware {
	return func(server *echo.Echo) {
		server.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
			Skipper: func(e echo.Context) bool {
				return strings.Contains(e.Path(), "ping")
			},
			CustomTimeFormat: "2006-01-02T15:04:05.1483386-00:00",
			Format: `{ "time":"${time_custom}", "level" :"Info" ,"method":"${method}", "uri":"${uri}",` +
				fmt.Sprintf(`"service": %q }`,
					cfg.ProjectName) + "\n",
		}))
	}
}

func WithRecover() Middleware {
	return func(server *echo.Echo) {
		server.Use(echoMiddleware.Recover())
	}
}

func HTTPErrorHandler(err error, ctx echo.Context) {
	var apiError httpserver.RestErr
	switch value := err.(type) {
	case *echo.HTTPError:
		apiError = httpserver.NewRestError(value.Error(), value.Code, err.Error())
	case exceptions.DuplicatedException:
		apiError = httpserver.NewRestError(err.Error(), http.StatusConflict, "conflict")
	case httpserver.RestErr:
		apiError = value
	case exceptions.NotFoundException:
		apiError = httpserver.NewNotFoundError(err.Error())
	case exceptions.InternalServerException:
		apiError = httpserver.NewUnauthorizedError(err.Error())
	default:
		apiError = httpserver.NewInternalServerError(err.Error(), err)
	}

	ctx.JSON(apiError.Status(), apiError)
}
