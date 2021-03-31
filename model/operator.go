package model

import (
	"time"

	"github.com/google/uuid"
)

type Operator struct {
	ID        string `gorm:"primary_key;not null;index"`
	UserID    string `gorm:"not null;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewOperator(userID string) *Operator {
	return &Operator{
		ID:     uuid.New().String(),
		UserID: userID,
	}
}
