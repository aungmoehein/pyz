package dbm

import (
	"errors"
	"time"
)

// Time wraps time.Time for pretty JONS formatting
type Time struct {
	time.Time
}

// Scan implements the SQL Scanner interface for TransactionType
func (wt *Time) Scan(value interface{}) error {
	if t, ok := value.(time.Time); ok {
		wt.Time = t
		return nil
	}

	return errors.New("Unable to scan formatted time")
}
