package services_test

import (
	"errors"
	"testing"

	"github.com/sebastianreh/user-balance-api/internal/app/services"
	"github.com/sebastianreh/user-balance-api/internal/domain/report"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
	"github.com/sebastianreh/user-balance-api/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMigrationReportService_GenerateAndSendReport(t *testing.T) {
	log := logger.NewLogger()
	reportBody := "Total Records Processed: 5000\nTotal Users Updated: 200"

	t.Run("it sends the report successfully", func(t *testing.T) {
		emailServiceMock := mocks.NewEmailServiceMock()
		reportService := services.NewMigrationReportService(log, emailServiceMock)

		summary := report.MigrationSummary{
			TotalRecords: 5000,
			UsersUpdated: 200,
		}
		to := []string{"recipient@example.com"}

		emailServiceMock.On("SendEmail", to, "Migration Report", reportBody).
			Return(nil)

		err := reportService.GenerateAndSendReport(summary, to)

		assert.Nil(t, err)
		emailServiceMock.AssertCalled(t, "SendEmail", to, "Migration Report", reportBody)
	})

	t.Run("it returns error when email sending fails", func(t *testing.T) {
		emailServiceMock := mocks.NewEmailServiceMock()
		reportService := services.NewMigrationReportService(log, emailServiceMock)

		summary := report.MigrationSummary{
			TotalRecords: 5000,
			UsersUpdated: 200,
		}
		to := []string{"recipient@example.com"}
		expectedError := errors.New("failed to send email")

		emailServiceMock.On("SendEmail", to, "Migration Report", reportBody).
			Return(expectedError)

		err := reportService.GenerateAndSendReport(summary, to)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not send report email")
		emailServiceMock.AssertCalled(t, "SendEmail", to, "Migration Report", reportBody)
	})
}
