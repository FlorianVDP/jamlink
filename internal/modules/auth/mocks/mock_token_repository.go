package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	tokenDomain "jamlink-backend/internal/modules/auth/domain/token"
)

type MockTokenRepository struct {
	mock.Mock
}

func (m *MockTokenRepository) FindByToken(id string) (*tokenDomain.Token, error) {
	args := m.Called(id)
	foundUser := args.Get(0)
	if foundUser == nil {
		return nil, args.Error(1)
	}
	return foundUser.(*tokenDomain.Token), args.Error(1)
}

func (m *MockTokenRepository) Create(token *tokenDomain.Token) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockTokenRepository) DeleteByID(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockTokenRepository) DeleteUserTokens(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}
