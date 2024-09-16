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

func NewUserHandler(log logger.Logger, service services.UserService) *UserHandler {
	return &UserHandler{
		log:     log,
		service: service,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user by providing first name, last name, and email
// @Tags users
// @Accept json
// @Produce json
// @Param user body user.User true "User Request Body"
// @Success 201 {object} user.CreationResponse "User created successfully with the user_id"
// @Failure 400 {object} exceptions.BadRequestException "Invalid input"
// @Failure 500 {object} exceptions.InternalServerException "Internal server error"
// @Router /users/create [post]
func (u *UserHandler) CreateUser(ctx echo.Context) error {
	userEntity, err := validateUserRequest(ctx)
	if err != nil {
		u.log.ErrorAt(err, userHandlerName, "CreateUser")
		exception := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(exception.Code(), exception)
	}

	createdID, err := u.service.CreateUser(ctx.Request().Context(), userEntity)
	if err != nil {
		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception)
	}

	response := user.CreationResponse{
		UserID: createdID,
	}

	return ctx.JSON(http.StatusCreated, response)
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Updates user details such as first name, last name, and email
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body user.User true "User Request Body"
// @Success 200 "User updated successfully"
// @Failure 400 {object} exceptions.BadRequestException  "Invalid request or missing user ID"
// @Failure 404 {object} exceptions.NotFoundException "User not found"
// @Failure 500 {object} exceptions.InternalServerException"Internal server error"
// @Router /users/{id} [put]
func (u *UserHandler) UpdateUser(ctx echo.Context) error {
	id, err := validateUserIDRequest(ctx)
	if err != nil {
		u.log.ErrorAt(err, userHandlerName, "UpdateUser")
		exception := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(exception.Code(), exception)
	}

	userEntity, err := validateUserRequest(ctx)
	if err != nil {
		u.log.ErrorAt(err, userHandlerName, "UpdateUser")
		exception := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(exception.Code(), exception)
	}

	userEntity.ID = id
	err = u.service.UpdateUser(ctx.Request().Context(), userEntity)
	if err != nil {
		if strings.Contains(err.Error(), user.NotFoundError) {
			exception := exceptions.NewBadRequestException(err.Error())
			return ctx.JSON(exception.Code(), exception)
		}

		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception)
	}

	return ctx.NoContent(http.StatusOK)
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Retrieves a user's details by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} user.User "User details"
// @Failure 400 {object} exceptions.BadRequestException "Invalid request or missing user ID"
// @Failure 404 {object} exceptions.NotFoundException "User not found"
// @Failure 500 {object} exceptions.InternalServerException "Internal server error"
// @Router /users/{id} [get]
func (u *UserHandler) GetUser(ctx echo.Context) error {
	id, err := validateUserIDRequest(ctx)
	if err != nil {
		u.log.ErrorAt(err, userHandlerName, "GetUser")
		exception := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(exception.Code(), exception)
	}

	userEntity, err := u.service.GetUser(ctx.Request().Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), user.NotFoundError) {
			exception := exceptions.NewNotFoundException(err.Error())
			return ctx.JSON(exception.Code(), exception)
		}

		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception)
	}

	return ctx.JSON(http.StatusOK, userEntity)
}

// DeleteUser godoc
// @Summary Delete a user by ID
// @Description Soft delete a user by marking them as deleted
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 "No Content"
// @Failure 400 {object} exceptions.BadRequestException "Invalid request or missing user ID"
// @Failure 404 {object} exceptions.NotFoundException "User not found"
// @Failure 500 {object} exceptions.InternalServerException "Internal server error"
// @Router /users/{id} [delete]
func (u *UserHandler) DeleteUser(ctx echo.Context) error {
	id, err := validateUserIDRequest(ctx)
	if err != nil {
		u.log.ErrorAt(err, userHandlerName, "DeleteUser")
		exception := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(exception.Code(), exception)
	}

	err = u.service.DeleteUser(ctx.Request().Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), user.NotFoundError) {
			exception := exceptions.NewNotFoundException(err.Error())
			return ctx.JSON(exception.Code(), exception)
		}

		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception)
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
