package tokenInvariants

import "time"

func ValidateToken(expDate time.Time) error {
	if err := TokenValidationExpDate(expDate); err != nil {
		return err
	}

	return nil
}
