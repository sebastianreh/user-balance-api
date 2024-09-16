package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/sebastianreh/user-balance-api/internal/domain/user"

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

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Create a new transaction for a user with a specified amount and datetime
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body transaction.Transaction true "Transaction Request Body"
// @Success 201 "No Content"
// @Failure 400 {object} exceptions.BadRequestException "Invalid request or business rule violation"
// @Failure 409 {object} exceptions.DuplicatedException "Transaction already exists"
// @Failure 500 {object} exceptions.InternalServerException "Internal server error"
// @Router /transactions/create [post]
func (t *TransactionHandler) CreateTransaction(ctx echo.Context) error {
	transactionEntity, err := validateTransactionRequest(ctx)
	if err != nil {
		t.log.ErrorAt(err, transactionHandlerName, "CreateTransaction")
		exception := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(exception.Code(), exception)
	}

	err = t.service.CreateTransaction(ctx.Request().Context(), transactionEntity)
	if err != nil {
		if strings.Contains(err.Error(), transaction.ZeroAmountError) {
			exception := exceptions.NewBadRequestException(err.Error())
			return ctx.JSON(exception.Code(), exception)
		}

		if strings.Contains(err.Error(), transaction.DuplicateTransactionError) {
			exception := exceptions.NewDuplicatedException(err.Error())
			return ctx.JSON(exception.Code(), exception)
		}

		if strings.Contains(err.Error(), user.NotFoundError) {
			exception := exceptions.NewBadRequestException(err.Error())
			return ctx.JSON(exception.Code(), exception)
		}

		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception)
	}

	return ctx.NoContent(http.StatusCreated)
}

// UpdateTransaction godoc
// @Summary Update an existing transaction
// @Description Update an existing transaction by ID with new data such as amount and datetime
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Param transaction body transaction.Transaction true "Transaction Request Body"
// @Success 200 "No Content"
// @Failure 400 {object} exceptions.BadRequestException "Invalid request or business rule violation"
// @Failure 500 {object} exceptions.InternalServerException "Internal server error"
// @Router /transactions/{id} [put]
func (t *TransactionHandler) UpdateTransaction(ctx echo.Context) error {
	id, err := validateTransactionIDRequest(ctx)
	if err != nil {
		t.log.ErrorAt(err, transactionHandlerName, "UpdateTransaction")
		exception := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(exception.Code(), exception)
	}

	transactionEntity, err := validateTransactionRequest(ctx)
	if err != nil {
		t.log.ErrorAt(err, transactionHandlerName, "CreateTransaction")
		exception := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(exception.Code(), exception)
	}

	transactionEntity.ID = id
	err = t.service.UpdateTransaction(ctx.Request().Context(), transactionEntity)
	if err != nil {
		if strings.Contains(err.Error(), transaction.NotFoundError) ||
			strings.Contains(err.Error(), transaction.ZeroAmountError) {
			exception := exceptions.NewBadRequestException(err.Error())
			return ctx.JSON(exception.Code(), exception)
		}

		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception)
	}

	return ctx.NoContent(http.StatusOK)
}

// GetTransaction godoc
// @Summary Get a transaction by ID
// @Description Retrieve transaction details by its ID
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} transaction.Transaction "Transaction details"
// @Failure 400 {object} exceptions.BadRequestException "Invalid request or business rule violation"
// @Failure 404 {object} exceptions.NotFoundException "Transaction not found"
// @Failure 500 {object} exceptions.InternalServerException "Internal server error"
// @Router /transactions/{id} [get]
func (t *TransactionHandler) GetTransaction(ctx echo.Context) error {
	id, err := validateTransactionIDRequest(ctx)
	if err != nil {
		t.log.ErrorAt(err, transactionHandlerName, "GetTransaction")
		exception := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(exception.Code(), exception)
	}

	transactionEntity, err := t.service.GetTransaction(ctx.Request().Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), transaction.NotFoundError) {
			exception := exceptions.NewNotFoundException(err.Error())
			return ctx.JSON(exception.Code(), exception)
		}

		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception)
	}

	return ctx.JSON(http.StatusOK, transactionEntity)
}

// DeleteTransaction godoc
// @Summary Delete a transaction by ID
// @Description Soft delete a transaction by its ID, marking it as deleted
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 "No Content"
// @Failure 400 {object} exceptions.BadRequestException "Invalid request or business rule violation"
// @Failure 404 {object} exceptions.NotFoundException "Transaction not found"
// @Failure 500 {object} exceptions.InternalServerException "Internal server error"
// @Router /transactions/{id} [delete]
func (t *TransactionHandler) DeleteTransaction(ctx echo.Context) error {
	id, err := validateTransactionIDRequest(ctx)
	if err != nil {
		t.log.ErrorAt(err, transactionHandlerName, "DeleteTransaction")
		exception := exceptions.NewBadRequestException(err.Error())
		return ctx.JSON(exception.Code(), exception)
	}

	err = t.service.DeleteTransaction(ctx.Request().Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), transaction.NotFoundError) {
			exception := exceptions.NewNotFoundException(err.Error())
			return ctx.JSON(exception.Code(), exception)
		}

		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception)
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
