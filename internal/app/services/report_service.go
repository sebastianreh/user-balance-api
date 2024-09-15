package services

import (
	"fmt"
	"github.com/sebastianreh/user-balance-api/internal/domain/report"
	"github.com/sebastianreh/user-balance-api/pkg/email"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
	"strings"
)

const (
	MigrationReportServiceName = "MigrationReportService"
)

type MigrationReportService interface {
	GenerateAndSendReport(migrationSummary report.MigrationSummary, to []string) error
}

type migrationReportService struct {
	log          logger.Logger
	emailService email.EmailService
}

func NewMigrationReportService(log logger.Logger, emailService email.EmailService) MigrationReportService {
	return &migrationReportService{
		log:          log,
		emailService: emailService,
	}
}

func (s *migrationReportService) GenerateAndSendReport(summary report.MigrationSummary, to []string) error {
	subject := "Migration Report"
	body := s.generateReportBody(summary)
	err := s.emailService.SendEmail(to, subject, body)
	if err != nil {
		err := fmt.Errorf("could not send report email, error: %w", err)
		s.log.ErrorAt(err, MigrationReportServiceName, "GenerateAndSendReport")
		return err
	}

	return nil
}

func (s *migrationReportService) generateReportBody(summary report.MigrationSummary) string {
	reportEmailBody := []string{
		fmt.Sprintf("Total Records Processed: %d", summary.TotalRecords),
		fmt.Sprintf("Total Users Updated: %d", summary.UsersUpdated),
	}

	return strings.Join(reportEmailBody, "\n")
}
