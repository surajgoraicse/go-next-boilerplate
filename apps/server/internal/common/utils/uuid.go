package utils

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func UUIDToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	val, err := uuid.FromBytes(u.Bytes[:])
	if err != nil {
		return ""
	}
	return val.String()
}

func StringToUUID(s string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{
		Bytes: parsed,
		Valid: true,
	}, nil
}
