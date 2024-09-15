package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sebastianreh/user-balance-api/cmd/httpserver/exceptions"
	"github.com/sebastianreh/user-balance-api/internal/app/services"
	"github.com/sebastianreh/user-balance-api/internal/domain/user"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
	customStr "github.com/sebastianreh/user-balance-api/pkg/strings"
)

const (
	userHandlerName = "UserHandler"
)

type UserHandler struct {
	service services.UserService
	log     logger.Logger
}

func NewUserHandler(logger logger.Logger, service services.UserService) *UserHandler {
	return &UserHandler{
		log:     logger,
		service: service,
	}
}

func (u *UserHandler) CreateUser(ctx echo.Context) error {
	userEntity, err := validateUserRequest(ctx)
	if err != nil {
		u.log.ErrorAt(err, userHandlerName, "CreateUser")
		exception := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(exception.Code(), exception.Error())
	}

	createdID, err := u.service.CreateUser(ctx.Request().Context(), userEntity)
	if err != nil {
		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception.Error())
	}

	return ctx.JSON(http.StatusCreated, struct {
		ID string `json:"user_id"`
	}{createdID})
}

func (u *UserHandler) UpdateUser(ctx echo.Context) error {
	id, err := validateUserIDRequest(ctx)
	if err != nil {
		u.log.ErrorAt(err, userHandlerName, "UpdateUser")
		exception := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(exception.Code(), exception.Error())
	}

	userEntity, err := validateUserRequest(ctx)
	if err != nil {
		u.log.ErrorAt(err, userHandlerName, "UpdateUser")
		exception := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(exception.Code(), exception.Error())
	}

	userEntity.ID = id
	err = u.service.UpdateUser(ctx.Request().Context(), userEntity)
	if err != nil {
		if strings.Contains(err.Error(), user.NotFoundError) {
			exception := exceptions.NewBadRequestException(err.Error())
			return ctx.JSON(exception.Code(), exception.Error())
		}

		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

func (u *UserHandler) GetUser(ctx echo.Context) error {
	id, err := validateUserIDRequest(ctx)
	if err != nil {
		u.log.ErrorAt(err, userHandlerName, "GetUser")
		exception := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(exception.Code(), exception.Error())
	}

	userEntity, err := u.service.GetUser(ctx.Request().Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), user.NotFoundError) {
			exception := exceptions.NewNotFoundException(err.Error())
			return ctx.JSON(exception.Code(), exception.Error())
		}

		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception.Error())
	}

	return ctx.JSON(http.StatusOK, userEntity)
}

func (u *UserHandler) DeleteUser(ctx echo.Context) error {
	id, err := validateUserIDRequest(ctx)
	if err != nil {
		u.log.ErrorAt(err, userHandlerName, "DeleteUser")
		exception := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(exception.Code(), exception.Error())
	}

	err = u.service.DeleteUser(ctx.Request().Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), user.NotFoundError) {
			exception := exceptions.NewNotFoundException(err.Error())
			return ctx.JSON(exception.Code(), exception.Error())
		}

		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

func validateUserRequest(ctx echo.Context) (user.User, error) {
	var userEntity user.User
	if err := ctx.Bind(&userEntity); err != nil {
		return userEntity, echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if customStr.IsEmpty(userEntity.FirstName) {
		return userEntity, errors.New("first name is required")
	}

	if customStr.IsEmpty(userEntity.LastName) {
		return userEntity, errors.New("last name is required")
	}

	if customStr.IsEmpty(userEntity.Email) {
		return userEntity, errors.New("email is required")
	}

	return userEntity, nil
}

func validateUserIDRequest(ctx echo.Context) (string, error) {
	id := ctx.Param("id")

	if customStr.IsEmpty(id) {
		return id, errors.New("missing param id")
	}

	return id, nil
}
