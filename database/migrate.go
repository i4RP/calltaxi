package database

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/sckacr/calltaxi/model"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Passenger{},
		&model.Operator{},
		&model.Request{},
	).Error
}
