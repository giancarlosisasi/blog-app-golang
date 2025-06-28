package utils

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func FromStrToUUID(value string) (uuid.UUID, error) {
	uid, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, ErrInvalidUUID
	}
	return uid, nil
}

func FromStrToPGUUID(value string) (pgtype.UUID, error) {
	uid, err := FromStrToUUID(value)
	if err != nil {
		return pgtype.UUID{Valid: false}, err
	}

	pguid := pgtype.UUID{
		Bytes: uid,
		Valid: true,
	}

	return pguid, nil
}

func FromStrToPGText(value string) pgtype.Text {
	return pgtype.Text{
		String: value,
		Valid:  true,
	}
}
