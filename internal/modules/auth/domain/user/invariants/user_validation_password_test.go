package userInvariants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password    string
		expectedErr error
	}{
		// Valid password
		{"Abcdef1!", nil},
		{"StrongPass123$", nil},
		{"A@bCdeFg1234567890", nil},

		// Invalid cases
		{"", ErrInvalidUserPassword},
		{"short1!", ErrShortPassword},
		{"thispasswordiswaytoolongtobeacceptablebecauseitexceedssixtyfourcharacters1234567890", ErrLongPassword},
		{"nouppercase1!", ErrNoUpperCase},
		{"NOLOWERCASE1!", ErrNoLowerCase},
		{"NoDigits!", ErrNoDigit},
		{"NoSpecial123", ErrNoSpecialChar},
	}

	for _, tt := range tests {
		err := ValidatePassword(tt.password)
		if tt.expectedErr == nil {
			assert.NoError(t, err, "Expected no error for password: %s", tt.password)
		} else {
			assert.ErrorIs(t, err, tt.expectedErr, "Expected %v for password: %s", tt.expectedErr, tt.password)
		}
	}
}
