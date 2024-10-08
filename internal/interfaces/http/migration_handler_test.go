package http_test

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sebastianreh/user-balance-api/internal/domain/report"

	"github.com/labstack/echo/v4"
	"github.com/sebastianreh/user-balance-api/internal/app/services"
	localHttp "github.com/sebastianreh/user-balance-api/internal/interfaces/http"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
	"github.com/sebastianreh/user-balance-api/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMigrationHandler_UploadMigrationCSV(t *testing.T) {
	log := logger.NewLogger()

	t.Run("it processes the CSV successfully and sends the report to the default email", func(t *testing.T) {
		serviceMock := mocks.NewMigrationServiceMock()
		migrationServiceMock := mocks.NewReportServiceMock()
		migrationReport := report.MigrationSummary{
			TotalRecords: 1000, UsersUpdated: 200,
		}

		rec, ctx := createMultipartFile(t, "test.csv", "1,1,100,2023-09-14T20:00:00Z")

		serviceMock.On("ProcessBalance", mock.Anything, mock.Anything).Return(migrationReport, nil)
		migrationServiceMock.On("GenerateAndSendReport", migrationReport, mock.Anything).Return(nil)

		handler := localHttp.NewMigrationHandler(log, serviceMock, migrationServiceMock)
		err := handler.UploadMigrationCSV(ctx)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		serviceMock.AssertCalled(t, "ProcessBalance", mock.Anything, mock.Anything)
	})

	t.Run("it processes the CSV successfully and sends the report to email addresses specified"+
		" in the X-User-Emails header.", func(t *testing.T) {
		serviceMock := mocks.NewMigrationServiceMock()
		migrationServiceMock := mocks.NewReportServiceMock()
		migrationReport := report.MigrationSummary{
			TotalRecords: 1000, UsersUpdated: 200,
		}

		rec, ctx := createMultipartFile(t, "test.csv", "1,1,100,2023-09-14T20:00:00Z")
		ctx.Request().Header.Set("X-Destination-Emails", "test1@example.com,test2@example.com")

		serviceMock.On("ProcessBalance", mock.Anything, mock.Anything).Return(migrationReport, nil)
		migrationServiceMock.On("GenerateAndSendReport", migrationReport, mock.Anything).Return(nil)

		handler := localHttp.NewMigrationHandler(log, serviceMock, migrationServiceMock)
		err := handler.UploadMigrationCSV(ctx)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		serviceMock.AssertCalled(t, "ProcessBalance", mock.Anything, mock.Anything)
	})

	t.Run("it returns an error for an invalid file format", func(t *testing.T) {
		serviceMock := mocks.NewMigrationServiceMock()
		migrationServiceMock := mocks.NewReportServiceMock()

		rec, ctx := createMultipartFile(t, "test.txt", "1,1,100,2023-09-14T20:00:00Z")

		handler := localHttp.NewMigrationHandler(log, serviceMock, migrationServiceMock)
		err := handler.UploadMigrationCSV(ctx)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("it returns an error for an empty CSV", func(t *testing.T) {
		serviceMock := mocks.NewMigrationServiceMock()
		migrationServiceMock := mocks.NewReportServiceMock()

		rec, ctx := createMultipartFile(t, "test.csv", "")

		handler := localHttp.NewMigrationHandler(log, serviceMock, migrationServiceMock)
		err := handler.UploadMigrationCSV(ctx)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("it returns an error for a CSV with the wrong format", func(t *testing.T) {
		serviceMock := mocks.NewMigrationServiceMock()
		migrationServiceMock := mocks.NewReportServiceMock()

		rec, ctx := createMultipartFile(t, "test.csv", "1,1,100") // Missing one column

		handler := localHttp.NewMigrationHandler(log, serviceMock, migrationServiceMock)
		err := handler.UploadMigrationCSV(ctx)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("it returns an error when the service fails", func(t *testing.T) {
		serviceMock := mocks.NewMigrationServiceMock()
		migrationServiceMock := mocks.NewReportServiceMock()

		rec, ctx := createMultipartFile(t, "test.csv", "1,1,100,2023-09-14T20:00:00Z")

		expectedError := errors.New(services.ReadFileError)
		serviceMock.On("ProcessBalance", mock.Anything, mock.Anything).Return(
			report.MigrationSummary{}, expectedError)

		handler := localHttp.NewMigrationHandler(log, serviceMock, migrationServiceMock)
		err := handler.UploadMigrationCSV(ctx)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		serviceMock.AssertCalled(t, "ProcessBalance", mock.Anything, mock.Anything)
	})

	t.Run("it returns an internal server error when migrationReportService fails unexpectedly", func(t *testing.T) {
		serviceMock := mocks.NewMigrationServiceMock()
		migrationServiceMock := mocks.NewReportServiceMock()
		migrationReport := report.MigrationSummary{
			TotalRecords: 1000, UsersUpdated: 200,
		}

		rec, ctx := createMultipartFile(t, "test.csv", "1,1,100,2023-09-14T20:00:00Z")

		expectedError := errors.New("migration report service error")
		serviceMock.On("ProcessBalance", mock.Anything, mock.Anything).Return(migrationReport, nil)
		migrationServiceMock.On("GenerateAndSendReport", migrationReport, mock.Anything).Return(expectedError)

		handler := localHttp.NewMigrationHandler(log, serviceMock, migrationServiceMock)
		err := handler.UploadMigrationCSV(ctx)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		serviceMock.AssertCalled(t, "ProcessBalance", mock.Anything, mock.Anything)
	})
}

func createMultipartFile(t *testing.T, filename string, content string) (*httptest.ResponseRecorder, echo.Context) {
	e := echo.New()
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	assert.NoError(t, err)

	_, err = part.Write([]byte(content))
	assert.NoError(t, err)

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	return rec, ctx
}
