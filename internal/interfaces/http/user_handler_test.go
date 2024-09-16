package http_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/sebastianreh/user-balance-api/cmd/httpserver"
	"github.com/sebastianreh/user-balance-api/internal/domain/user"
	localHttp "github.com/sebastianreh/user-balance-api/internal/interfaces/http"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
	"github.com/sebastianreh/user-balance-api/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserHandler_CreateUser(t *testing.T) {
	log := logger.NewLogger()

	t.Run("it creates a new user successfully", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()

		userRequest := user.User{
			FirstName: "user",
			LastName:  "lastname",
			Email:     "user@example.com",
		}

		requestBytes, _ := json.Marshal(userRequest)
		context, rec := httpserver.SetupAsRecorder(http.MethodPost, "/users/create", "", string(requestBytes))
		serviceMock.On("CreateUser", mock.Anything, userRequest).Return("1", nil)

		handler := localHttp.NewUserHandler(log, serviceMock)
		err := handler.CreateUser(context)

		var response struct {
			ID string `json:"user_id"`
		}
		_ = json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Nil(t, err)
		assert.Equal(t, "1", response.ID)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("it returns bad request when validation fails", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()

		invalidUserRequest := `{"FirstName": "user", "LastName": ""}` // Invalid request
		context, rec := httpserver.SetupAsRecorder(http.MethodPost, "/users/create", "", invalidUserRequest)

		handler := localHttp.NewUserHandler(log, serviceMock)
		err := handler.CreateUser(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		serviceMock.AssertNotCalled(t, "CreateUser", mock.Anything, mock.Anything)
	})

	t.Run("it returns internal server error when service fails", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()

		userRequest := user.User{
			FirstName: "user",
			LastName:  "lastname",
			Email:     "user@example.com",
		}

		requestBytes, _ := json.Marshal(userRequest)
		context, rec := httpserver.SetupAsRecorder(http.MethodPost, "/users/create", "", string(requestBytes))
		serviceMock.On("CreateUser", mock.Anything, userRequest).Return("", errors.New("service failure"))

		handler := localHttp.NewUserHandler(log, serviceMock)
		err := handler.CreateUser(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestUserHandler_GetUser(t *testing.T) {
	log := logger.NewLogger()

	t.Run("it gets user successfully", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()

		userID := "1"
		expectedResponse := user.User{
			FirstName: "user",
			LastName:  "lastname",
			Email:     "user@example.com",
		}

		context, rec := httpserver.SetupAsRecorder(http.MethodGet, "/:id", userID, "")
		serviceMock.On("GetUser", mock.Anything, userID).Return(expectedResponse, nil)

		handler := localHttp.NewUserHandler(log, serviceMock)
		err := handler.GetUser(context)

		var response user.User
		_ = json.Unmarshal(rec.Body.Bytes(), &response)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedResponse, response)
	})

	t.Run("it returns bad request when user ID validation fails", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()

		context, rec := httpserver.SetupAsRecorder(http.MethodGet, "/:id", "", "") // Missing user ID

		handler := localHttp.NewUserHandler(log, serviceMock)
		err := handler.GetUser(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		serviceMock.AssertNotCalled(t, "GetUser", mock.Anything, mock.Anything)
	})

	t.Run("it returns not found when user is not found", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()

		userID := "1"
		expectedError := errors.New(user.NotFoundError)

		context, rec := httpserver.SetupAsRecorder(http.MethodGet, "/:id", userID, "")
		serviceMock.On("GetUser", mock.Anything, userID).Return(user.User{}, expectedError)

		handler := localHttp.NewUserHandler(log, serviceMock)
		err := handler.GetUser(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("it returns internal server error when service fails", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()

		userID := "1"
		expectedError := errors.New("service failure")

		context, rec := httpserver.SetupAsRecorder(http.MethodGet, "/:id", userID, "")
		serviceMock.On("GetUser", mock.Anything, userID).Return(user.User{}, expectedError)

		handler := localHttp.NewUserHandler(log, serviceMock)
		err := handler.GetUser(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestUserHandler_UpdateUser(t *testing.T) {
	log := logger.NewLogger()

	t.Run("it updates a user successfully", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()

		userRequest := user.User{
			ID:        "1",
			FirstName: "user",
			LastName:  "lastname",
			Email:     "user@example.com",
		}

		requestBytes, _ := json.Marshal(userRequest)
		context, rec := httpserver.SetupAsRecorder(http.MethodPut, "/:id", userRequest.ID, string(requestBytes))
		serviceMock.On("UpdateUser", mock.Anything, userRequest).Return(nil)

		handler := localHttp.NewUserHandler(log, serviceMock)
		err := handler.UpdateUser(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("it returns bad request when user ID validation fails", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()

		userRequest := user.User{
			FirstName: "user",
			LastName:  "lastname",
			Email:     "user@example.com",
		}

		requestBytes, _ := json.Marshal(userRequest)
		context, rec := httpserver.SetupAsRecorder(http.MethodPut, "/:id", "", string(requestBytes)) // Missing ID
		handler := localHttp.NewUserHandler(log, serviceMock)

		err := handler.UpdateUser(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		serviceMock.AssertNotCalled(t, "UpdateUser", mock.Anything, mock.Anything)
	})

	t.Run("it returns bad request when request validation fails", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()

		context, rec := httpserver.SetupAsRecorder(http.MethodPut, "/:id", "1", "{}")
		handler := localHttp.NewUserHandler(log, serviceMock)

		err := handler.UpdateUser(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		serviceMock.AssertNotCalled(t, "UpdateUser", mock.Anything, mock.Anything)
	})

	t.Run("it returns not found when user is not found", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()

		userRequest := user.User{
			ID:        "1",
			FirstName: "user",
			LastName:  "lastname",
			Email:     "user@example.com",
		}

		requestBytes, _ := json.Marshal(userRequest)
		context, rec := httpserver.SetupAsRecorder(http.MethodPut, "/:id", userRequest.ID, string(requestBytes))
		expectedError := errors.New(user.NotFoundError)
		serviceMock.On("UpdateUser", mock.Anything, userRequest).Return(expectedError)

		handler := localHttp.NewUserHandler(log, serviceMock)
		err := handler.UpdateUser(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("it returns internal server error when service fails", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()

		userRequest := user.User{
			ID:        "1",
			FirstName: "user",
			LastName:  "lastname",
			Email:     "user@example.com",
		}

		requestBytes, _ := json.Marshal(userRequest)
		context, rec := httpserver.SetupAsRecorder(http.MethodPut, "/:id", userRequest.ID, string(requestBytes))
		expectedError := errors.New("service failure")
		serviceMock.On("UpdateUser", mock.Anything, userRequest).Return(expectedError)

		handler := localHttp.NewUserHandler(log, serviceMock)
		err := handler.UpdateUser(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestUserHandler_DeleteUser(t *testing.T) {
	log := logger.NewLogger()

	t.Run("it deletes user successfully", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()

		userID := "1"
		context, rec := httpserver.SetupAsRecorder(http.MethodDelete, "/:id", userID, "")
		serviceMock.On("DeleteUser", mock.Anything, userID).Return(nil)

		handler := localHttp.NewUserHandler(log, serviceMock)
		err := handler.DeleteUser(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("it returns bad request when user ID validation fails", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()

		// Invalid or missing user ID
		context, rec := httpserver.SetupAsRecorder(http.MethodDelete, "/:id", "", "")
		handler := localHttp.NewUserHandler(log, serviceMock)

		err := handler.DeleteUser(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		serviceMock.AssertNotCalled(t, "DeleteUser", mock.Anything, mock.Anything)
	})

	t.Run("it returns not found when user is not found", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()

		userID := "1"
		context, rec := httpserver.SetupAsRecorder(http.MethodDelete, "/:id", userID, "")
		expectedError := errors.New(user.NotFoundError)
		serviceMock.On("DeleteUser", mock.Anything, userID).Return(expectedError)

		handler := localHttp.NewUserHandler(log, serviceMock)
		err := handler.DeleteUser(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("it returns internal server error when service fails", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()

		userID := "1"
		context, rec := httpserver.SetupAsRecorder(http.MethodDelete, "/:id", userID, "")
		expectedError := errors.New("service failure")
		serviceMock.On("DeleteUser", mock.Anything, userID).Return(expectedError)

		handler := localHttp.NewUserHandler(log, serviceMock)
		err := handler.DeleteUser(context)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
