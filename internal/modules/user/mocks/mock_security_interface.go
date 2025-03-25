package mocks

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockSecurityService struct {
	mock.Mock
}

func (m *MockSecurityService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockSecurityService) CheckPassword(password, hash string) bool {
	args := m.Called(password, hash)
	return args.Bool(0)
}

func (m *MockSecurityService) GenerateJWT(id *uuid.UUID, email *string, duration time.Duration, tokenType string) (string, error) {
	args := m.Called(id, email, duration, tokenType)
	return args.String(0), args.Error(1)
}

func (m *MockSecurityService) ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(jwt.MapClaims), args.Error(1)
}

func (m *MockSecurityService) GetJWTInfo(token string) (uuid.UUID, error) {
	args := m.Called(token)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockSecurityService) GenerateSecureRandomString(n int) (string, error) {
	args := m.Called(n)
	return args.Get(0).(string), args.Error(1)
}
