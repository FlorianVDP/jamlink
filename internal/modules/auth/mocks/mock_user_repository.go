package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	userDomain "jamlink-backend/internal/modules/auth/domain/user"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(email string) (*userDomain.User, error) {
	args := m.Called(email)
	foundUser := args.Get(0)
	if foundUser == nil {
		return nil, args.Error(1)
	}
	return foundUser.(*userDomain.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uuid.UUID) (*userDomain.User, error) {
	args := m.Called(id)
	foundUser := args.Get(0)
	if foundUser == nil {
		return nil, args.Error(1)
	}
	return foundUser.(*userDomain.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *userDomain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *userDomain.User) error {
	args := m.Called(user)
	return args.Error(0)
}
