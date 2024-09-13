package http

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sebastianreh/user-balance-api/cmd/httpserver/exceptions"
	"github.com/sebastianreh/user-balance-api/internal/app/services"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
	"github.com/sebastianreh/user-balance-api/pkg/strings"
	"net/http"
	"time"
)

const (
	balanceHandlerName = "BalanceHandler"
)

type BalanceHandler struct {
	service services.BalanceService
	log     logger.Logger
}

func NewBalanceHandler(logger logger.Logger, service services.BalanceService) *BalanceHandler {
	return &BalanceHandler{
		log:     logger,
		service: service,
	}
}

func (h *BalanceHandler) GetUserBalanceWithOptions(ctx echo.Context) error {
	if isWithOptionsRequest(ctx) {
		return h.HandleGetUserBalanceWithOptions(ctx)
	}

	return h.HandleGetUserBalanceWithoutOptions(ctx)
}

func (h *BalanceHandler) HandleGetUserBalanceWithOptions(ctx echo.Context) error {
	id, fromDate, toDate, err := validateBalanceWithOptionsRequest(ctx)
	if err != nil {
		err := exceptions.NewBadRequestException(err.Error())
		h.log.ErrorAt(err, balanceHandlerName, "HandleGetUserBalanceWithoutOptions")
		return ctx.JSON(err.Code(), err.Error())
	}

	balance, err := h.service.GetBalanceByUserIDWithOptions(ctx.Request().Context(), id, fromDate, toDate)
	if err != nil {
		if err.Error() == services.UserNotFound {
			err := exceptions.NewNotFoundException(err.Error())
			return ctx.JSON(err.Code(), err.Error())
		}

		err := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(err.Code(), err.Error())
	}

	return ctx.JSON(http.StatusOK, balance)
}

func (h *BalanceHandler) HandleGetUserBalanceWithoutOptions(ctx echo.Context) error {
	id, err := validateUserBalanceRequest(ctx)
	if err != nil {
		err := exceptions.NewBadRequestException(err.Error())
		h.log.ErrorAt(err, balanceHandlerName, "HandleGetUserBalanceWithoutOptions")
		return ctx.JSON(err.Code(), err.Error())
	}

	balance, err := h.service.GetBalanceByUserID(ctx.Request().Context(), id)
	if err != nil {
		if err.Error() == services.UserNotFound {
			err := exceptions.NewNotFoundException(err.Error())
			return ctx.JSON(err.Code(), err.Error())
		}

		err := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(err.Code(), err.Error())
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

	if strings.IsEmpty(id) {
		return id, fromDate, toDate, errors.New("missing param id")
	}

	if strings.IsEmpty(fromDate) || strings.IsEmpty(toDate) {
		return id, fromDate, toDate, errors.New("missing date values")
	}

	err = validateDates(fromDate, toDate)
	if err != nil {
		return id, fromDate, toDate, err
	}

	return id, fromDate, toDate, nil
}

func validateDates(fromDate, toDate string) error {
	layout := "2006-01-02T15:04:05Z"

	fromTime, err := time.Parse(layout, fromDate)
	if err != nil {
		return fmt.Errorf("invalid fromDate format: %v", err)
	}

	toTime, err := time.Parse(layout, toDate)
	if err != nil {
		return fmt.Errorf("invalid toDate format: %v", err)
	}

	if fromTime.After(toTime) {
		return fmt.Errorf("fromDate cannot be after toDate")
	}

	return nil
}

func validateUserBalanceRequest(ctx echo.Context) (string, error) {
	id := ctx.Param("user_id")

	if strings.IsEmpty(id) {
		return id, errors.New("missing param id")
	}

	return id, nil
}
