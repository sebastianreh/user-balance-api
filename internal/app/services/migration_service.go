package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/sebastianreh/user-balance-api/internal/domain/transaction"
	"github.com/sebastianreh/user-balance-api/internal/domain/user"
	"github.com/sebastianreh/user-balance-api/pkg/csv"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
	"io"
	"mime/multipart"
)

const (
	migrationServiceName = "MigrationService"
	ReadFileError        = "error reading file"
)

type MigrationService interface {
	ProcessBalance(ctx context.Context, file *multipart.FileHeader) error
}

type migrationService struct {
	log                   logger.Logger
	userRepository        user.Repository
	transactionRepository transaction.Repository
	csvProcessor          csv.CsvProcessor
}

func NewMigrationService(logger logger.Logger, userRepository user.Repository,
	transactionRepository transaction.Repository, csvProcessor csv.CsvProcessor) MigrationService {
	return &migrationService{
		log:                   logger,
		userRepository:        userRepository,
		transactionRepository: transactionRepository,
		csvProcessor:          csvProcessor,
	}
}

func (s *migrationService) ProcessBalance(ctx context.Context, file *multipart.FileHeader) error {
	records, err := s.readFile(file)
	if err != nil {
		s.log.ErrorAt(fmt.Errorf("error reading file: %s", err.Error()),
			migrationServiceName, "ProcessBalance")
		return errors.New(ReadFileError)
	}

	for _, record := range records {
		err = s.userRepository.Save(ctx, user.CreateUserByRecord(record))
		if err != nil {
			return err
		}

		userTransaction, err := transaction.CreateTransactionByRecord(record)
		if err != nil {
			s.log.ErrorAt(fmt.Errorf("error creating transaction by record: %s", err.Error()),
				migrationServiceName, "ProcessBalance")
		}

		err = s.transactionRepository.Save(ctx, userTransaction)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *migrationService) readFile(file *multipart.FileHeader) ([][]string, error) {
	records := make([][]string, 0)

	src, err := file.Open()
	if err != nil {
		return records, err
	}
	defer src.Close()

	var fileBytes []byte
	fileBytes, err = io.ReadAll(src)
	if err != nil {
		return records, err
	}

	records, err = s.csvProcessor.CsvBytesToRecords(fileBytes, recordValidator)
	if err != nil {
		return records, err
	}

	return records, nil
}
