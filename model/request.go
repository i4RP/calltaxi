package model

import (
	"time"

	"github.com/google/uuid"
)

type Request struct {
	ID          string `gorm:"primary_key;not null;index"`
	PassengerID string `gorm:"not null;"`
	OperatorID  string
	Address     string
	Finished    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewRequest() *Request {
	return &Request{
		ID: uuid.New().String(),
	}
}
