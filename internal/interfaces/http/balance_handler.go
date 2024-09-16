package http

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sebastianreh/user-balance-api/cmd/httpserver/exceptions"
	"github.com/sebastianreh/user-balance-api/internal/app/services"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
	customStr "github.com/sebastianreh/user-balance-api/pkg/strings"
)

const (
	balanceHandlerName   = "BalanceHandler"
	TimeLayoutUTC        = "2006-01-02T15:04:05Z"
	TimeLayoutWithOffset = "2006-01-02T15:04:05-07:00"
)

type BalanceHandler struct {
	service services.BalanceService
	log     logger.Logger
}

func NewBalanceHandler(log logger.Logger, service services.BalanceService) *BalanceHandler {
	return &BalanceHandler{
		log:     log,
		service: service,
	}
}

// GetUserBalanceWithOptions godoc
// @Summary Get user balance with optional date filters
// @Description Get the balance of a user.
// If "from" and "to" query parameters are provided, the balance is filtered by the specified date range.
// @Tags balances
// @Param user_id path string true "User ID"
// @Param from query string false "Start date in ISO8601 format (YYYY-MM-DDThh:mm:ssZ)"
// @Param to query string false "End date in ISO8601 format (YYYY-MM-DDThh:mm:ssZ)"
// @Success 200 {object} balance.UserBalance
// @Failure 400 {object} exceptions.BadRequestException
// @Failure 404 {object} exceptions.NotFoundException
// @Failure 500 {object} exceptions.InternalServerException
// @Router /users/{user_id}/balance [get]
func (h *BalanceHandler) GetUserBalanceWithOptions(ctx echo.Context) error {
	if isWithOptionsRequest(ctx) {
		return h.HandleGetUserBalanceWithOptions(ctx)
	}

	return h.HandleGetUserBalanceWithoutOptions(ctx)
}

func (h *BalanceHandler) HandleGetUserBalanceWithOptions(ctx echo.Context) error {
	id, fromDate, toDate, err := validateBalanceWithOptionsRequest(ctx)
	if err != nil {
		exception := exceptions.NewBadRequestException(err.Error())
		h.log.ErrorAt(exception, balanceHandlerName, "HandleGetUserBalanceWithoutOptions")
		return ctx.JSON(exception.Code(), exception)
	}

	balance, err := h.service.GetBalanceByUserIDWithOptions(ctx.Request().Context(), id, fromDate, toDate)
	if err != nil {
		if err.Error() == services.UserNotFound {
			exception := exceptions.NewBadRequestException(err.Error())
			return ctx.JSON(exception.Code(), exception)
		}

		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception)
	}

	return ctx.JSON(http.StatusOK, balance)
}

func (h *BalanceHandler) HandleGetUserBalanceWithoutOptions(ctx echo.Context) error {
	id, err := validateUserBalanceRequest(ctx)
	if err != nil {
		exception := exceptions.NewBadRequestException(err.Error())
		h.log.ErrorAt(exception, balanceHandlerName, "HandleGetUserBalanceWithoutOptions")
		return ctx.JSON(exception.Code(), exception)
	}

	balance, err := h.service.GetBalanceByUserID(ctx.Request().Context(), id)
	if err != nil {
		if err.Error() == services.UserNotFound {
			exception := exceptions.NewNotFoundException(err.Error())
			return ctx.JSON(exception.Code(), exception)
		}

		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception)
	}

	return ctx.JSON(http.StatusOK, balance)
}

func isWithOptionsRequest(ctx echo.Context) bool {
	fromDate := ctx.QueryParam("from")
	ToDate := ctx.QueryParam("to")

	if fromDate == "" && ToDate == "" {
		return false
	}

	return true
}

func validateBalanceWithOptionsRequest(ctx echo.Context) (id, fromDate, toDate string, err error) {
	id = ctx.Param("user_id")
	fromDate = ctx.QueryParam("from")
	toDate = ctx.QueryParam("to")

	if customStr.IsEmpty(id) {
		return id, fromDate, toDate, errors.New("missing param id")
	}

	if customStr.IsEmpty(fromDate) || customStr.IsEmpty(toDate) {
		return id, fromDate, toDate, errors.New("missing date values")
	}

	err = validateDates(fromDate, toDate)
	if err != nil {
		return id, fromDate, toDate, err
	}

	return id, fromDate, toDate, nil
}

func validateDates(fromDate, toDate string) error {
	fromTime, err := time.Parse(TimeLayoutUTC, fromDate)
	if err != nil {
		fromTime, err = time.Parse(TimeLayoutWithOffset, fromDate)
		if err != nil {
			return fmt.Errorf("invalid fromDate format: %v", err)
		}
	}

	toTime, err := time.Parse(TimeLayoutUTC, toDate)
	if err != nil {
		toTime, err = time.Parse(TimeLayoutWithOffset, toDate)
		if err != nil {
			return fmt.Errorf("invalid toDate format: %v", err)
		}
	}

	if fromTime.After(toTime) {
		return fmt.Errorf("fromDate cannot be after toDate")
	}

	return nil
}

func validateUserBalanceRequest(ctx echo.Context) (string, error) {
	id := ctx.Param("user_id")

	if customStr.IsEmpty(id) {
		return id, errors.New("missing param id")
	}

	return id, nil
}
