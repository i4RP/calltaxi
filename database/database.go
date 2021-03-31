package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func New(url string) (*gorm.DB, error) {
	return gorm.Open("postgres", url)
}
