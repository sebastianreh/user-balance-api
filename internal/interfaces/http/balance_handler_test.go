package http_test

import (
	"encoding/json"
	"errors"
	"github.com/sebastianreh/user-balance-api/cmd/httpserver"
	"github.com/sebastianreh/user-balance-api/internal/app/services"
	"github.com/sebastianreh/user-balance-api/internal/domain/balance"
	localHttp "github.com/sebastianreh/user-balance-api/internal/interfaces/http"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
	"github.com/sebastianreh/user-balance-api/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func TestBalanceHandler_GetUserBalanceWithoutOptions(t *testing.T) {
	log := logger.NewLogger()

	t.Run("it gets user balance without options successfully", func(t *testing.T) {
		serviceMock := mocks.NewBalanceServiceMock()
		userID := "1"
		expectedBalance := balance.UserBalance{Balance: 100.00}

		context, rec := httpserver.SetupAsRecorderWithIdField(http.MethodGet,
			"/balances", userID, "", "user_id")
		serviceMock.On("GetBalanceByUserID", mock.Anything, userID).Return(expectedBalance, nil)

		handler := localHttp.NewBalanceHandler(log, serviceMock)
		err := handler.GetUserBalanceWithOptions(context)

		var response balance.UserBalance
		_ = json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedBalance, response)
	})

	t.Run("it returns bad request for missing user ID", func(t *testing.T) {
		serviceMock := mocks.NewBalanceServiceMock()

		context, rec := httpserver.SetupAsRecorderWithIdField(http.MethodGet, "/balances", "", "", "user_id")
		handler := localHttp.NewBalanceHandler(log, serviceMock)

		err := handler.GetUserBalanceWithOptions(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("it returns not found when user balance is not found", func(t *testing.T) {
		serviceMock := mocks.NewBalanceServiceMock()
		userID := "1"

		context, rec := httpserver.SetupAsRecorderWithIdField(http.MethodGet, "/balances", userID, "", "user_id")
		serviceMock.On("GetBalanceByUserID", mock.Anything, userID).Return(balance.UserBalance{}, errors.New(services.UserNotFound))

		handler := localHttp.NewBalanceHandler(log, serviceMock)
		err := handler.GetUserBalanceWithOptions(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("it returns internal server error when service fails", func(t *testing.T) {
		serviceMock := mocks.NewBalanceServiceMock()
		userID := "1"

		context, rec := httpserver.SetupAsRecorderWithIdField(http.MethodGet, "/balances", userID, "", "user_id")
		serviceMock.On("GetBalanceByUserID", mock.Anything, userID).Return(balance.UserBalance{}, errors.New("service error"))

		handler := localHttp.NewBalanceHandler(log, serviceMock)
		err := handler.GetUserBalanceWithOptions(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestBalanceHandler_GetUserBalanceWithOptions(t *testing.T) {
	log := logger.NewLogger()
	userID := "1"

	t.Run("it gets user balance with options successfully with UTC dates", func(t *testing.T) {
		serviceMock := mocks.NewBalanceServiceMock()
		userID := "1"
		expectedBalance := balance.UserBalance{Balance: 100.00}

		fromDate := "2024-05-02T15:04:05-03:00"
		toDate := "2024-09-02T20:13:28-03:00"
		queryParams := map[string]string{
			"from": fromDate,
			"to":   toDate,
		}

		context, rec := httpserver.SetupAsRecorderWithDynamicQueryParams(http.MethodGet, "/balances", userID, queryParams, "")
		serviceMock.On("GetBalanceByUserIDWithOptions", mock.Anything, userID, fromDate, toDate).Return(expectedBalance, nil)

		handler := localHttp.NewBalanceHandler(log, serviceMock)
		err := handler.GetUserBalanceWithOptions(context)

		var response balance.UserBalance
		_ = json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedBalance, response)
	})

	t.Run("it gets user balance with options successfully with Timezone Dates", func(t *testing.T) {
		serviceMock := mocks.NewBalanceServiceMock()
		expectedBalance := balance.UserBalance{Balance: 100.00}

		fromDate := "2024-05-02T15:04:05-03:00"
		toDate := "2024-09-02T20:13:28-03:00"
		queryParams := map[string]string{
			"from": fromDate,
			"to":   toDate,
		}

		context, rec := httpserver.SetupAsRecorderWithDynamicQueryParams(http.MethodGet, "/balances", userID, queryParams, "")
		serviceMock.On("GetBalanceByUserIDWithOptions", mock.Anything, userID, fromDate, toDate).Return(expectedBalance, nil)

		handler := localHttp.NewBalanceHandler(log, serviceMock)
		err := handler.GetUserBalanceWithOptions(context)

		var response balance.UserBalance
		_ = json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedBalance, response)
	})

	t.Run("it returns bad request for missing from date value", func(t *testing.T) {
		serviceMock := mocks.NewBalanceServiceMock()

		fromDate := "2024-05-02T15:04:05-03:00"
		queryParams := map[string]string{
			"from": fromDate,
			"to":   "",
		}

		context, rec := httpserver.SetupAsRecorderWithDynamicQueryParams(http.MethodGet, "/balances", userID, queryParams, "")
		handler := localHttp.NewBalanceHandler(log, serviceMock)

		err := handler.GetUserBalanceWithOptions(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("it returns bad request for missing to date value", func(t *testing.T) {
		serviceMock := mocks.NewBalanceServiceMock()

		fromDate := "2024-05-02T15:04:05-03:00"
		queryParams := map[string]string{
			"from": fromDate,
			"to":   "",
		}

		context, rec := httpserver.SetupAsRecorderWithDynamicQueryParams(http.MethodGet, "/balances", userID, queryParams, "")
		handler := localHttp.NewBalanceHandler(log, serviceMock)

		err := handler.GetUserBalanceWithOptions(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("it returns bad request for invalid date range", func(t *testing.T) {
		serviceMock := mocks.NewBalanceServiceMock()

		fromDate := "2024-05-02T15:04:05-03:00"
		toDate := "2023-09-02T20:13:28-03:00"
		queryParams := map[string]string{
			"from": fromDate,
			"to":   toDate,
		}

		context, rec := httpserver.SetupAsRecorderWithDynamicQueryParams(http.MethodGet, "/balances", userID, queryParams, "")
		handler := localHttp.NewBalanceHandler(log, serviceMock)

		err := handler.GetUserBalanceWithOptions(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("it returns internal server error when service fails in GetUserBalanceWithOptions", func(t *testing.T) {
		serviceMock := mocks.NewBalanceServiceMock()
		fromDate := "2024-05-02T15:04:05Z"
		toDate := "2024-09-02T20:13:28Z"
		queryParams := map[string]string{
			"from": fromDate,
			"to":   toDate,
		}

		context, rec := httpserver.SetupAsRecorderWithDynamicQueryParams(http.MethodGet, "/balances", userID, queryParams, "")
		serviceMock.On("GetBalanceByUserIDWithOptions", mock.Anything, userID, fromDate, toDate).Return(balance.UserBalance{}, errors.New("service error"))

		handler := localHttp.NewBalanceHandler(log, serviceMock)
		err := handler.GetUserBalanceWithOptions(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
