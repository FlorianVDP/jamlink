package tokenInvariants

import (
	"errors"
	"time"
)

var (
	ErrTokenExpDateBeforeCurrentDate = errors.New("token expire date is before current date")
)

func TokenValidationExpDate(expDate time.Time) error {
	if expDate.Before(time.Now()) {
		return ErrTokenExpDateBeforeCurrentDate
	}

	return nil
}
