package mocks

import (
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

type CsvProcessorMock struct {
	mock.Mock
}

func NewCsvProcessorMock() *CsvProcessorMock {
	return new(CsvProcessorMock)
}

func (m *CsvProcessorMock) ReadFile(file *multipart.FileHeader, recordValidator func(record []string) error) ([][]string, error) {
	args := m.Called(file, recordValidator)
	return args.Get(0).([][]string), args.Error(1)
}
