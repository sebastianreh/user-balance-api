package http

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sebastianreh/user-balance-api/cmd/httpserver/exceptions"
	"github.com/sebastianreh/user-balance-api/internal/app/services"
	"github.com/sebastianreh/user-balance-api/internal/domain/transaction"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
)

const (
	migrationHandlerName = "MigrationHandler"
	fileColumns          = 4
)

type MigrationHandler struct {
	log           logger.Logger
	service       services.MigrationService
	reportService services.MigrationReportService
}

func NewMigrationHandler(log logger.Logger, service services.MigrationService,
	reportService services.MigrationReportService) *MigrationHandler {
	return &MigrationHandler{
		log:           log,
		service:       service,
		reportService: reportService,
	}
}

// UploadMigrationCSV handles the uploading of a CSV file containing migration data. godoc
// It reads the file, processes the migration, and sends a migration report to the specified email addresses.
//
// @Summary      Upload Migration CSV
// @Description  This endpoint allows uploading a CSV file that contains migration data.
//               The system processes the CSV file, migrates the necessary data, and sends a report
//               to the email addresses specified in the "X-Destination-Emails" header.
// @Tags         Migration
// @Accept       multipart/form-data
// @Produce      application/json
// @Param        file         formData   file   true  "CSV file with migration data"
// @Param        X-User-Emails  header    string true  "Comma-separated list of email addresses to send the migration report"
// @Success      200 "No content"
// @Failure      400 {object}  exceptions.BadRequestException {message=string} "Bad request (e.g., invalid CSV file format)"
// @Failure      500 {object}  exceptions.InternalServerException {message=string} "Internal server error"

func (h *MigrationHandler) UploadMigrationCSV(ctx echo.Context) error {
	file, err := validateFile(ctx)
	if err != nil {
		exception := exceptions.NewBadRequestException(err.Error())
		h.log.ErrorAt(exception, migrationHandlerName, "UploadMigrationsCSV")
		return ctx.JSON(exception.Code(), exception)
	}

	migrationReport, err := h.service.ProcessBalance(ctx.Request().Context(), file)
	if err != nil {
		if err.Error() == services.ReadFileError || strings.Contains(err.Error(), transaction.DuplicateTransactionError) {
			exception := exceptions.NewBadRequestException(err.Error())
			return ctx.JSON(exception.Code(), exception)
		}

		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception)
	}

	reportDestinations := getDestinationEmailsFromRequestHeader(ctx)
	err = h.reportService.GenerateAndSendReport(migrationReport, reportDestinations)
	if err != nil {
		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception)
	}

	return ctx.NoContent(http.StatusOK)
}

func validateFile(ctx echo.Context) (*multipart.FileHeader, error) {
	file, err := ctx.FormFile("file")
	if err != nil {
		return nil, err
	}

	if !strings.Contains(file.Filename, ".csv") {
		return nil, errors.New("the file must be csv")
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	csvReader := csv.NewReader(src)
	record, err := csvReader.Read()
	if err == io.EOF {
		return nil, errors.New("empty file")
	}

	if err != nil {
		return nil, fmt.Errorf("cannot read file - error: %s", err.Error())
	}

	if len(record) != fileColumns {
		return nil, errors.New("the file is not in the correct format")
	}

	return file, nil
}

func getDestinationEmailsFromRequestHeader(ctx echo.Context) []string {
	emailsHeader := ctx.Request().Header.Get("X-Destination-Emails")
	emails := strings.Split(emailsHeader, ",")
	for i, email := range emails {
		emails[i] = strings.TrimSpace(email)
	}

	return emails
}
