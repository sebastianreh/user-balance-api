package mocks

import (
	"github.com/stretchr/testify/mock"
)

type EmailServiceMock struct {
	mock.Mock
}

func NewEmailServiceMock() *EmailServiceMock {
	return new(EmailServiceMock)
}

func (m *EmailServiceMock) SendEmail(to []string, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}
