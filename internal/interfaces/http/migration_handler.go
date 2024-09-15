package http

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sebastianreh/user-balance-api/cmd/httpserver/exceptions"
	"github.com/sebastianreh/user-balance-api/internal/app/services"
	"github.com/sebastianreh/user-balance-api/internal/domain/transaction"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
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

func NewMigrationHandler(logger logger.Logger, service services.MigrationService,
	reportService services.MigrationReportService) *MigrationHandler {
	return &MigrationHandler{
		log:           logger,
		service:       service,
		reportService: reportService,
	}
}

func (h *MigrationHandler) UploadMigrationCSV(ctx echo.Context) error {
	file, err := validateFile(ctx)
	if err != nil {
		exception := exceptions.NewBadRequestException(err.Error())
		h.log.ErrorAt(exception, migrationHandlerName, "UploadMigrationsCSV")
		return ctx.JSON(exception.Code(), exception.Error())
	}

	migrationReport, err := h.service.ProcessBalance(ctx.Request().Context(), file)
	if err != nil {
		if err.Error() == services.ReadFileError || strings.Contains(err.Error(), transaction.DuplicateTransactionError) {
			exception := exceptions.NewBadRequestException(err.Error())
			return ctx.JSON(exception.Code(), exception.Error())
		}

		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception.Error())
	}

	err = h.reportService.GenerateAndSendReport(migrationReport, []string{})
	if err != nil {
		exception := exceptions.NewInternalServerException(err.Error())
		return ctx.JSON(exception.Code(), exception.Error())
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
