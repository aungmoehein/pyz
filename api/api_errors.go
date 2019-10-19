package api

import (
	"fmt"
	"net/http"

	"hackathon.com/pyz/dbm"
)

// Error specify error code / http status status etc
type Error struct {
	status  int
	code    string
	message string
	details interface{}
}

// Wraps another error inside APIError
func (e *Error) Wraps(details interface{}) Error {
	switch v := details.(type) {
	case dbm.Error:
		e.message = v.Error()
	default:
		e.details = v
	}

	return *e
}

// String return error description
func (e Error) Error() string {
	if e.details != nil {
		return fmt.Sprint(e.details)
	}

	return e.message
}

// ErrUnexpected is raised on unexpected errors
var ErrUnexpected = Error{
	code:    "ErrUnexpected",
	status:  http.StatusInternalServerError,
	message: "Internal server error occured",
}

// ErrURLParams is raised when no UID param in URL
var ErrURLParams = Error{
	code:    "ErrURLParams",
	status:  http.StatusBadRequest,
	message: "Unable to parse URL params",
}

// ErrEmptyRequestBody is raised when request JSON body is empty
var ErrEmptyRequestBody = Error{
	code:    "ErrEmptyRequestBody",
	status:  http.StatusBadRequest,
	message: "Request body is empty (no params)",
}

// ErrRequestParams is raised when parsing JSON request body params
// or when the params are not valid
var ErrRequestParams = Error{
	code:    "ErrRequestParams",
	status:  http.StatusBadRequest,
	message: "Unable to parse JSON params (Either mising or invalid)",
}

// ErrDatabase is raised on database errors
var ErrDatabase = Error{
	code:    "ErrDatabase",
	status:  http.StatusInternalServerError,
	message: "Database error occured",
}

// ErrDatabaseTx is raised when there are problems with database transaction
var ErrDatabaseTx = Error{
	code:    "ErrDatabaseTx",
	status:  http.StatusInternalServerError,
	message: "Unable to start / commit / rollback a database transaction",
}

// ErrInvalidUID is raised when parsing URL params fails
var ErrInvalidUID = Error{
	code:    "ErrInvalidUID",
	status:  http.StatusBadRequest,
	message: "Unable to parse UID param",
}
