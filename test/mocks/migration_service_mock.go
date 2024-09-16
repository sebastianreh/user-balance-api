package mocks

import (
	"context"
	"mime/multipart"

	"github.com/sebastianreh/user-balance-api/internal/domain/report"
	"github.com/stretchr/testify/mock"
)

type MigrationServiceMock struct {
	mock.Mock
}

func NewMigrationServiceMock() *MigrationServiceMock {
	return new(MigrationServiceMock)
}

func (m *MigrationServiceMock) ProcessBalance(ctx context.Context, file *multipart.FileHeader) (report.MigrationSummary, error) {
	args := m.Called(ctx, file)
	return args.Get(0).(report.MigrationSummary), args.Error(1)
}
