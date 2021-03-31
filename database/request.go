package database

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/sckacr/calltaxi/model"
)

type Request interface {
	FindByID(string) (*model.Request, error)
	FindLatest(string) (*model.Request, error)
	Create(*model.Request) (*model.Request, error)
	Update(*model.Request) (*model.Request, error)
	Delete(string) error
}

type RequestImpl struct {
	*gorm.DB
}

func NewRequest(db *gorm.DB) Request {
	return &RequestImpl{db}
}

func (r *RequestImpl) FindByID(id string) (*model.Request, error) {
	request := new(model.Request)

	err := r.Where("id = ?", id).Find(request).Error

	return request, err
}

func (r *RequestImpl) FindLatest(passengerID string) (*model.Request, error) {
	request := new(model.Request)

	err := r.Where("passenger_id = ?", passengerID).Order("created_at desc").Limit(1).Find(request).Error

	return request, err
}

func (r *RequestImpl) Create(request *model.Request) (*model.Request, error) {
	tx := r.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Create(request).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return request, tx.Commit().Error
}

func (r *RequestImpl) Update(request *model.Request) (*model.Request, error) {
	tx := r.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Where("id = ?", request.ID).Save(request).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return request, tx.Commit().Error
}

func (r *RequestImpl) Delete(id string) error {
	tx := r.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("id = ?", id).Delete(&model.Request{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
