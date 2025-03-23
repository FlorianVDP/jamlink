package mocks

import (
	"github.com/stretchr/testify/mock"
	"jamlink-backend/internal/modules/user/domain"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(email string) (*userDomain.User, error) {
	args := m.Called(email)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*userDomain.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *userDomain.User) error {
	args := m.Called(user)
	return args.Error(0)
}
