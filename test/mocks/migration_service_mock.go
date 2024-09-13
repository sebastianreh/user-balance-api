package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"mime/multipart"
)

type MigrationServiceMock struct {
	mock.Mock
}

func NewMigrationServiceMock() *MigrationServiceMock {
	return new(MigrationServiceMock)
}

func (m *MigrationServiceMock) ProcessBalance(ctx context.Context, file *multipart.FileHeader) error {
	args := m.Called(ctx, file)
	return args.Error(0)
}
