package dbm

import (
	"database/sql"
	"errors"
)

// Error Wraps error
type Error error

// ErrNoRowAffected is raised where update statement with no changes
var ErrNoRowAffected Error = sql.ErrNoRows

// ErrNotEnoughBalance is raised when insert / update with invalid valid
var ErrNotEnoughBalance Error = errors.New(
	"The wallet operation failed for not having enough balance",
)
