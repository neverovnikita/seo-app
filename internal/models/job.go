package models

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	Id        uuid.UUID `json:"id" db:"id"`
	EntityId  uuid.UUID `json:"entity_id" db:"entity_id"`
	Stage     string    `json:"stage" db:"stage"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
