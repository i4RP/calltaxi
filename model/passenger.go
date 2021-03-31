package model

import (
	"time"

	"github.com/google/uuid"
)

type Passenger struct {
	ID           string `gorm:"primary_key;not null;index"`
	UserID       string `gorm:"not null;"`
	Name         string
	ChangingName bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewPassenger(userID string) *Passenger {
	return &Passenger{
		ID:     uuid.New().String(),
		UserID: userID,
	}
}
