package util

import (
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetUUIDFromPGTypeUUID(id pgtype.UUID) uuid.UUID {
	ret, err := uuid.FromBytes(id.Bytes[:])

	if err != nil {
		log.Fatalln("Impossible error in UUID conversion")
	}

	return ret
}

func GetPGTypeUUIDFromString(id string) pgtype.UUID {
	ret, err := uuid.Parse(id)

	if err != nil {
		return pgtype.UUID{}
	}

	return pgtype.UUID{Bytes: ret, Valid: true}
}
