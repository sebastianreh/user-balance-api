package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sebastianreh/user-balance-api/cmd/httpserver/exceptions"
	"github.com/sebastianreh/user-balance-api/internal/app/services"
	"github.com/sebastianreh/user-balance-api/internal/domain/transaction"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
	customStr "github.com/sebastianreh/user-balance-api/pkg/strings"
)

const (
	transactionHandlerName = "TransactionHandler"
)

type TransactionHandler struct {
	service services.TransactionService
	log     logger.Logger
}

func NewTransactionHandler(log logger.Logger, service services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		log:     log,
		service: service,
	}
}

func (t *TransactionHandler) CreateTransaction(ctx echo.Context) error {
	transactionEntity, err := validateTransactionRequest(ctx)
	if err != nil {
		t.log.ErrorAt(err, transactionHandlerName, "CreateTransaction")
		exception := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(exception.Code(), exception.Error())
	}

	err = t.service.CreateTransaction(ctx.Request().Context(), transactionEntity)
	if err != nil {
		if strings.Contains(err.Error(), transaction.DuplicateTransactionError) ||
			strings.Contains(err.Error(), transaction.ZeroAmountError) {
			exception := exceptions.NewBadRequestException(err.Error())
			return ctx.JSON(exception.Code(), exception.Error())
		}

		t.log.ErrorAt(err, transactionHandlerName, "CreateTransaction")
		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception.Error())
	}

	return ctx.NoContent(http.StatusCreated)
}

func (t *TransactionHandler) UpdateTransaction(ctx echo.Context) error {
	id, err := validateTransactionIDRequest(ctx)
	if err != nil {
		t.log.ErrorAt(err, transactionHandlerName, "UpdateTransaction")
		exception := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(exception.Code(), exception.Error())
	}

	transactionEntity, err := validateTransactionRequest(ctx)
	if err != nil {
		t.log.ErrorAt(err, transactionHandlerName, "CreateTransaction")
		exception := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(exception.Code(), exception.Error())
	}

	transactionEntity.ID = id
	err = t.service.UpdateTransaction(ctx.Request().Context(), transactionEntity)
	if err != nil {
		if strings.Contains(err.Error(), transaction.NotFoundError) ||
			strings.Contains(err.Error(), transaction.ZeroAmountError) {
			exception := exceptions.NewBadRequestException(err.Error())
			return ctx.JSON(exception.Code(), exception.Error())
		}

		t.log.ErrorAt(err, transactionHandlerName, "UpdateTransaction")
		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

func (t *TransactionHandler) GetTransaction(ctx echo.Context) error {
	id, err := validateTransactionIDRequest(ctx)
	if err != nil {
		t.log.ErrorAt(err, transactionHandlerName, "GetTransaction")
		exception := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(exception.Code(), exception.Error())
	}

	transactionEntity, err := t.service.GetTransaction(ctx.Request().Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), transaction.NotFoundError) {
			exception := exceptions.NewNotFoundException(err.Error())
			return ctx.JSON(exception.Code(), exception.Error())
		}

		t.log.ErrorAt(err, transactionHandlerName, "GetTransaction")
		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception.Error())
	}

	return ctx.JSON(http.StatusOK, transactionEntity)
}

func (t *TransactionHandler) DeleteTransaction(ctx echo.Context) error {
	id, err := validateTransactionIDRequest(ctx)
	if err != nil {
		t.log.ErrorAt(err, transactionHandlerName, "DeleteTransaction")
		exception := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(exception.Code(), exception.Error())
	}

	err = t.service.DeleteTransaction(ctx.Request().Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), transaction.NotFoundError) {
			exception := exceptions.NewNotFoundException(err.Error())
			return ctx.JSON(exception.Code(), exception.Error())
		}

		t.log.ErrorAt(err, transactionHandlerName, "DeleteTransaction")
		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

func validateTransactionRequest(ctx echo.Context) (transaction.Transaction, error) {
	var transactionEntity transaction.Transaction
	if err := ctx.Bind(&transactionEntity); err != nil {
		return transactionEntity, echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if customStr.IsEmpty(transactionEntity.UserID) {
		return transactionEntity, errors.New("user ID is required")
	}

	if transactionEntity.Amount == 0 {
		return transactionEntity, errors.New("amount must be greater than zero")
	}

	if transactionEntity.DateTime.IsZero() {
		return transactionEntity, errors.New("date_time is required")
	}

	return transactionEntity, nil
}

func validateTransactionIDRequest(ctx echo.Context) (string, error) {
	id := ctx.Param("id")

	if customStr.IsEmpty(id) {
		return id, errors.New("missing param id")
	}

	return id, nil
}
