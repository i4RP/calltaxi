package database

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/sckacr/calltaxi/model"
)

type Passenger interface {
	FindByUserID(string) (*model.Passenger, error)
	Create(*model.Passenger) (*model.Passenger, error)
	Update(*model.Passenger) (*model.Passenger, error)
	Delete(string) error
}

type PassengerImpl struct {
	*gorm.DB
}

func NewPassenger(db *gorm.DB) Passenger {
	return &PassengerImpl{db}
}

func (p *PassengerImpl) FindByUserID(userID string) (*model.Passenger, error) {
	passenger := new(model.Passenger)

	err := p.Where("user_id = ?", userID).Find(passenger).Error

	return passenger, err
}

func (p *PassengerImpl) Create(passenger *model.Passenger) (*model.Passenger, error) {
	tx := p.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Create(passenger).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return passenger, tx.Commit().Error
}

func (p *PassengerImpl) Update(passenger *model.Passenger) (*model.Passenger, error) {
	tx := p.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Where("id = ?", passenger.ID).Save(passenger).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return passenger, tx.Commit().Error
}

func (p *PassengerImpl) Delete(id string) error {
	tx := p.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("id = ?", id).Delete(&model.Passenger{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
