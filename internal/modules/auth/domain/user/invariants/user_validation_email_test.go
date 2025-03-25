package userInvariants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email       string
		expectError bool
	}{
		// Valid cases
		{"test@example.com", false},
		{"user.name+tag@sub.domain.com", false},
		{"name@domain.co", false},
		{"a@b.io", false},

		// Invalid cases
		{"", true},
		{"plainaddress", true},
		{"@no-local-part.com", true},
		{"user@.nodomain", true},
		{"user@domain..com", true},
		{"user@domain", true},
		{"user@domain.c", true}, // TLD trop court
	}

	for _, tt := range tests {
		err := ValidateEmail(tt.email)
		if tt.expectError {
			assert.Error(t, err, "Expected error for email: %s", tt.email)
		} else {
			assert.NoError(t, err, "Expected no error for email: %s", tt.email)
		}
	}
}
