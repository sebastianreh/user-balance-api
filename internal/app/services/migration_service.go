package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"sync"

	"github.com/sebastianreh/user-balance-api/internal/domain/transaction"
	"github.com/sebastianreh/user-balance-api/internal/domain/user"
	"github.com/sebastianreh/user-balance-api/internal/infrastructure/config"
	"github.com/sebastianreh/user-balance-api/pkg/csv"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
)

const (
	migrationServiceName = "MigrationService"
	ReadFileError        = "error reading file"
)

type MigrationService interface {
	ProcessBalance(ctx context.Context, file *multipart.FileHeader) error
}

type migrationService struct {
	config                config.Config
	log                   logger.Logger
	userRepository        user.Repository
	transactionRepository transaction.Repository
	csvProcessor          csv.CsvProcessor
}

func NewMigrationService(config config.Config, logger logger.Logger, userRepository user.Repository,
	transactionRepository transaction.Repository, csvProcessor csv.CsvProcessor) MigrationService {
	return &migrationService{
		config:                config,
		log:                   logger,
		userRepository:        userRepository,
		transactionRepository: transactionRepository,
		csvProcessor:          csvProcessor,
	}
}

func (s *migrationService) ProcessBalance(ctx context.Context, file *multipart.FileHeader) error {
	records, err := s.readFile(file)
	if err != nil {
		s.log.ErrorAt(fmt.Errorf("error reading file: %s", err.Error()), migrationServiceName, "ProcessBalance")
		return errors.New(ReadFileError)
	}

	recordChan := make(chan []string)
	errChan := make(chan error, s.config.Workers.MigrationWorkersSize)

	var wg sync.WaitGroup

	for i := 0; i < s.config.Workers.MigrationWorkersSize; i++ {
		wg.Add(1)
		go s.processRecords(ctx, recordChan, errChan, &wg)
	}

	go func() {
		for _, record := range records {
			recordChan <- record
		}
		close(recordChan)
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	err = collectErrors(errChan)
	if err != nil {
		fmt.Println("Collected errors:", err)
		return err
	}

	return nil
}

func collectErrors(errChan <-chan error) error {
	var errorMessages []string
	for err := range errChan {
		if err != nil {
			errorMessages = append(errorMessages, err.Error())
		}
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf("errors found in balance process job - %s", strings.Join(errorMessages, ", "))
	}

	return nil
}

func (s *migrationService) processRecords(ctx context.Context, recordChan <-chan []string, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for record := range recordChan {
		err := s.userRepository.Save(ctx, user.CreateUserByRecord(record))
		if err != nil {
			errChan <- fmt.Errorf("error saving user: %w", err)
			return
		}

		userTransaction, err := transaction.CreateTransactionByRecord(record)
		if err != nil {
			s.log.ErrorAt(fmt.Errorf("error creating transaction by record: %s", err.Error()), migrationServiceName, "ProcessBalance")
			errChan <- err
			return
		}

		err = s.transactionRepository.Save(ctx, userTransaction)
		if err != nil {
			errChan <- fmt.Errorf("error saving transaction: %w", err)
			return
		}
	}

	errChan <- nil
}

/*
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

*/

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
