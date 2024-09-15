package middlewares

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/sebastianreh/user-balance-api/cmd/httpserver"
	"github.com/sebastianreh/user-balance-api/cmd/httpserver/exceptions"
	"net/http"
)

type Middleware func(echo *echo.Echo)

// AddMiddlewares build the middlewares of the server
func AddMiddlewares(server *echo.Echo, middlewares ...Middleware) {
	for _, middleware := range middlewares {
		middleware(server)
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
