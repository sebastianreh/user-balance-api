package http_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/sebastianreh/user-balance-api/cmd/httpserver"
	"github.com/sebastianreh/user-balance-api/internal/domain/transaction"
	localHttp "github.com/sebastianreh/user-balance-api/internal/interfaces/http"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
	"github.com/sebastianreh/user-balance-api/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionHandler_CreateTransaction(t *testing.T) {
	log := logger.NewLogger()
	now := time.Now().Truncate(0)

	t.Run("it creates a new transaction successfully", func(t *testing.T) {
		serviceMock := mocks.NewTransactionServiceMock()
		transactionRequest := transaction.Transaction{
			UserID:   "1",
			Amount:   100.00,
			DateTime: &now,
		}

		requestBytes, _ := json.Marshal(transactionRequest)
		context, rec := httpserver.SetupAsRecorder(http.MethodPost, "/transactions/create", "", string(requestBytes))
		serviceMock.On("CreateTransaction", mock.Anything, transactionRequest).Return(nil)

		handler := localHttp.NewTransactionHandler(log, serviceMock)
		err := handler.CreateTransaction(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("it returns bad request for invalid request body", func(t *testing.T) {
		serviceMock := mocks.NewTransactionServiceMock()

		context, rec := httpserver.SetupAsRecorder(http.MethodPost, "/transactions/create", "", "invalid body")
		handler := localHttp.NewTransactionHandler(log, serviceMock)

		err := handler.CreateTransaction(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("it returns bad request for missing user ID", func(t *testing.T) {
		serviceMock := mocks.NewTransactionServiceMock()

		transactionRequest := transaction.Transaction{
			Amount:   100.00,
			DateTime: &now,
		}

		requestBytes, _ := json.Marshal(transactionRequest)
		context, rec := httpserver.SetupAsRecorder(http.MethodPost, "/transactions/create", "", string(requestBytes))
		handler := localHttp.NewTransactionHandler(log, serviceMock)

		err := handler.CreateTransaction(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("it returns internal server error when service fails", func(t *testing.T) {
		serviceMock := mocks.NewTransactionServiceMock()

		transactionRequest := transaction.Transaction{
			UserID:   "1",
			Amount:   100.00,
			DateTime: &now,
		}

		requestBytes, _ := json.Marshal(transactionRequest)
		context, rec := httpserver.SetupAsRecorder(http.MethodPost, "/transactions/create", "", string(requestBytes))
		expectedError := errors.New("service error")
		serviceMock.On("CreateTransaction", mock.Anything, transactionRequest).Return(expectedError)

		handler := localHttp.NewTransactionHandler(log, serviceMock)
		err := handler.CreateTransaction(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestTransactionHandler_UpdateTransaction(t *testing.T) {
	log := logger.NewLogger()
	now := time.Now().Truncate(0)

	t.Run("it updates a transaction successfully", func(t *testing.T) {
		serviceMock := mocks.NewTransactionServiceMock()
		transactionRequest := transaction.Transaction{
			ID:       "1",
			UserID:   "1",
			Amount:   100.00,
			DateTime: &now,
		}

		requestBytes, _ := json.Marshal(transactionRequest)
		context, rec := httpserver.SetupAsRecorder(http.MethodPut, "/transactions/:id", transactionRequest.ID, string(requestBytes))
		serviceMock.On("UpdateTransaction", mock.Anything, transactionRequest).Return(nil)

		handler := localHttp.NewTransactionHandler(log, serviceMock)
		err := handler.UpdateTransaction(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("it returns bad request for invalid request body", func(t *testing.T) {
		serviceMock := mocks.NewTransactionServiceMock()

		context, rec := httpserver.SetupAsRecorder(http.MethodPut, "/transactions/:id", "1", "invalid body")
		handler := localHttp.NewTransactionHandler(log, serviceMock)

		err := handler.UpdateTransaction(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("it returns bad request for missing user ID", func(t *testing.T) {
		serviceMock := mocks.NewTransactionServiceMock()

		transactionRequest := transaction.Transaction{
			ID:       "1",
			Amount:   100.00,
			DateTime: &now,
		}

		requestBytes, _ := json.Marshal(transactionRequest)
		context, rec := httpserver.SetupAsRecorder(http.MethodPut, "/transactions/:id", transactionRequest.ID, string(requestBytes))
		handler := localHttp.NewTransactionHandler(log, serviceMock)

		err := handler.UpdateTransaction(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("it returns internal server error when service fails", func(t *testing.T) {
		serviceMock := mocks.NewTransactionServiceMock()

		transactionRequest := transaction.Transaction{
			ID:       "1",
			UserID:   "1",
			Amount:   100.00,
			DateTime: &now,
		}

		requestBytes, _ := json.Marshal(transactionRequest)
		context, rec := httpserver.SetupAsRecorder(http.MethodPut, "/transactions/:id", transactionRequest.ID, string(requestBytes))
		expectedError := errors.New("service error")
		serviceMock.On("UpdateTransaction", mock.Anything, transactionRequest).Return(expectedError)

		handler := localHttp.NewTransactionHandler(log, serviceMock)
		err := handler.UpdateTransaction(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestTransactionHandler_GetTransaction(t *testing.T) {
	log := logger.NewLogger()

	t.Run("it gets transaction successfully", func(t *testing.T) {
		serviceMock := mocks.NewTransactionServiceMock()
		now := time.Now().Truncate(0)

		transactionResponse := transaction.Transaction{
			ID:       "1",
			UserID:   "1",
			Amount:   100.00,
			DateTime: &now,
		}

		context, rec := httpserver.SetupAsRecorder(http.MethodGet, "/transactions/:id", "1", "")
		serviceMock.On("GetTransaction", mock.Anything, "1").Return(transactionResponse, nil)

		handler := localHttp.NewTransactionHandler(log, serviceMock)
		err := handler.GetTransaction(context)

		var response transaction.Transaction
		_ = json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, transactionResponse, response)
	})

	t.Run("it returns bad request for missing transaction ID", func(t *testing.T) {
		serviceMock := mocks.NewTransactionServiceMock()

		context, rec := httpserver.SetupAsRecorder(http.MethodGet, "/transactions/:id", "", "")
		handler := localHttp.NewTransactionHandler(log, serviceMock)

		err := handler.GetTransaction(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("it returns not found when transaction is not found", func(t *testing.T) {
		serviceMock := mocks.NewTransactionServiceMock()

		context, rec := httpserver.SetupAsRecorder(http.MethodGet, "/transactions/:id", "1", "")
		expectedError := errors.New(transaction.NotFoundError)
		serviceMock.On("GetTransaction", mock.Anything, "1").Return(transaction.Transaction{}, expectedError)

		handler := localHttp.NewTransactionHandler(log, serviceMock)
		err := handler.GetTransaction(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("it returns internal server error when service fails", func(t *testing.T) {
		serviceMock := mocks.NewTransactionServiceMock()

		context, rec := httpserver.SetupAsRecorder(http.MethodGet, "/transactions/:id", "1", "")
		expectedError := errors.New("service error")
		serviceMock.On("GetTransaction", mock.Anything, "1").Return(transaction.Transaction{}, expectedError)

		handler := localHttp.NewTransactionHandler(log, serviceMock)
		err := handler.GetTransaction(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestTransactionHandler_DeleteTransaction(t *testing.T) {
	log := logger.NewLogger()

	t.Run("it deletes transaction successfully", func(t *testing.T) {
		serviceMock := mocks.NewTransactionServiceMock()

		context, rec := httpserver.SetupAsRecorder(http.MethodDelete, "/transactions/:id", "1", "")
		serviceMock.On("DeleteTransaction", mock.Anything, "1").Return(nil)

		handler := localHttp.NewTransactionHandler(log, serviceMock)
		err := handler.DeleteTransaction(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("it returns bad request for missing transaction ID", func(t *testing.T) {
		serviceMock := mocks.NewTransactionServiceMock()

		context, rec := httpserver.SetupAsRecorder(http.MethodDelete, "/transactions/:id", "", "")
		handler := localHttp.NewTransactionHandler(log, serviceMock)

		err := handler.DeleteTransaction(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("it returns not found when transaction is not found", func(t *testing.T) {
		serviceMock := mocks.NewTransactionServiceMock()

		context, rec := httpserver.SetupAsRecorder(http.MethodDelete, "/transactions/:id", "1", "")
		expectedError := errors.New(transaction.NotFoundError)
		serviceMock.On("DeleteTransaction", mock.Anything, "1").Return(expectedError)

		handler := localHttp.NewTransactionHandler(log, serviceMock)
		err := handler.DeleteTransaction(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("it returns internal server error when service fails", func(t *testing.T) {
		serviceMock := mocks.NewTransactionServiceMock()

		context, rec := httpserver.SetupAsRecorder(http.MethodDelete, "/transactions/:id", "1", "")
		expectedError := errors.New("service error")
		serviceMock.On("DeleteTransaction", mock.Anything, "1").Return(expectedError)

		handler := localHttp.NewTransactionHandler(log, serviceMock)
		err := handler.DeleteTransaction(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
