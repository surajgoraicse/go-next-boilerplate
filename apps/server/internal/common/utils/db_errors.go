package utils

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// DbErrIsUniqueViolation (ErrUniqueViolation) checks if the error is a "unique violation" error.
func DbErrIsUniqueViolation(err error) bool {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		return pgErr.Code == "23505"
	}
	return false
}

// DbErrIsNotFound (ErrNoRows) checks if the error is a "not found" error.
func DbErrIsNotFound(err error) bool {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		return pgErr.Code == "23505"
	}
	return false
}
