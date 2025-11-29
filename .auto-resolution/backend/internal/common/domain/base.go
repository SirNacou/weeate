package domain

import (
	"time"

	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type Base struct {
	UUID
	Audit
	SoftDeleteableModel
}

func NewBase() Base {
	return Base{
		UUID: NewUUID(),
	}
}

type UUID struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuidv7()"`
}

func NewUUID() UUID {
	return UUID{
		ID: uuid.Must(uuid.NewV7()),
	}
}

type Audit struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SoftDeleteableModel struct {
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
