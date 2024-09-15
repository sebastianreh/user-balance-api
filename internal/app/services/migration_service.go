package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/sebastianreh/user-balance-api/internal/domain/report"
	"github.com/sebastianreh/user-balance-api/pkg/csv"
	"github.com/sebastianreh/user-balance-api/pkg/email"
	"mime/multipart"
	"strings"
	"sync"

	"github.com/sebastianreh/user-balance-api/internal/domain/transaction"
	"github.com/sebastianreh/user-balance-api/internal/domain/user"
	"github.com/sebastianreh/user-balance-api/internal/infrastructure/config"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
)

const (
	migrationServiceName = "MigrationService"
	ReadFileError        = "error reading file"
)

type MigrationService interface {
	ProcessBalance(ctx context.Context, file *multipart.FileHeader) (report.MigrationSummary, error)
}

type migrationService struct {
	config                config.Config
	log                   logger.Logger
	userRepository        user.Repository
	transactionRepository transaction.Repository
	csvProcessor          csv.CsvProcessor
	emailService          email.EmailService
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

func (s *migrationService) ProcessBalance(ctx context.Context, file *multipart.FileHeader) (report.MigrationSummary, error) {
	var migrationSummary report.MigrationSummary
	records, err := s.csvProcessor.ReadFile(file, recordValidator)
	if err != nil {
		s.log.ErrorAt(fmt.Errorf("error reading file: %s", err.Error()), migrationServiceName, "ProcessBalance")
		return migrationSummary, errors.New(ReadFileError)
	}

	batches := createBatches(records, s.config.Workers.MigrationWorkerBatchSize)

	errChan := make(chan error, len(batches))
	processedChan := make(chan map[string][]transaction.Transaction, len(batches))
	var wg sync.WaitGroup

	for _, batch := range batches {
		wg.Add(1)
		go func(batch [][]string) {
			defer wg.Done()
			s.processBatch(ctx, batch, processedChan, errChan)
		}(batch)
	}

	go func() {
		wg.Wait()
		close(processedChan)
		close(errChan)
	}()

	if err = collectErrors(errChan); err != nil {
		s.log.ErrorAt(fmt.Errorf("error during batch processing: %s", err.Error()), migrationServiceName, "ProcessBalance")
		return migrationSummary, err
	}

	migrationSummary = generateReport(processedChan)

	return migrationSummary, nil
}

func (s *migrationService) processBatch(ctx context.Context, batch [][]string,
	processedChan chan<- map[string][]transaction.Transaction, errChan chan<- error) {
	var transactions []transaction.Transaction
	processedTransactions := make(map[string][]transaction.Transaction)

	for _, record := range batch {
		userTransaction, err := transaction.CreateTransactionByRecord(record)
		if err != nil {
			errChan <- fmt.Errorf("error creating transaction by record: %w", err)
			return
		}
		transactions = append(transactions, userTransaction)
		processedTransactions[userTransaction.UserID] = append(processedTransactions[userTransaction.UserID], userTransaction)
	}

	err := s.transactionRepository.SaveBatch(ctx, transactions)
	if err != nil {
		errChan <- fmt.Errorf("error saving transaction batch: %w", err)
	}
	processedChan <- processedTransactions

	errChan <- nil
}

func createBatches(records [][]string, batchSize int) [][][]string {
	var batches [][][]string
	for batchSize < len(records) {
		records, batches = records[batchSize:], append(batches, records[0:batchSize:batchSize])
	}

	batches = append(batches, records)
	return batches
}

func collectErrors(errChan <-chan error) error {
	var errorMessages []string
	for err := range errChan {
		if err != nil {
			errorMessages = append(errorMessages, err.Error())
		}
	}

	errorSet := make(map[string]bool)
	var uniqueErrors []string

	for _, errMsg := range errorMessages {
		if _, exists := errorSet[errMsg]; !exists {
			errorSet[errMsg] = true
			uniqueErrors = append(uniqueErrors, errMsg)
		}
	}

	if len(uniqueErrors) > 0 {
		return errors.New(strings.Join(uniqueErrors, ", "))
	}

	return nil
}

func generateReport(processedChan <-chan map[string][]transaction.Transaction) report.MigrationSummary {
	var summary report.MigrationSummary

	// Created a map to track unique users since it's using a batch processor
	uniqueUsers := make(map[string]bool)
	for reportMap := range processedChan {
		for userID, userTransactions := range reportMap {
			summary.TotalRecords += len(userTransactions)

			// Track unique users
			if !uniqueUsers[userID] {
				uniqueUsers[userID] = true
				summary.UsersUpdated++
			}
		}
	}

	return summary
}
