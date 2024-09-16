package mocks

import (
	"github.com/sebastianreh/user-balance-api/internal/domain/report"
	"github.com/stretchr/testify/mock"
)

type ReportServiceMock struct {
	mock.Mock
}

func NewReportServiceMock() *ReportServiceMock {
	return new(ReportServiceMock)
}

func (m *ReportServiceMock) GenerateAndSendReport(summary report.MigrationSummary, to []string) error {
	args := m.Called(summary, to)
	return args.Error(0)
}
