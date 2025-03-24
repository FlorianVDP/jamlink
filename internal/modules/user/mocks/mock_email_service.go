package mocks

import (
	"github.com/stretchr/testify/mock"
	"jamlink-backend/internal/shared/email"
)

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) Send(to string, template email.TemplateType, lang string, data map[string]string) error {
	args := m.Called(to, template, lang, data)
	return args.Error(0)
}
