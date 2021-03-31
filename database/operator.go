package database

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/sckacr/calltaxi/model"
)

type Operator interface {
	FindAll() ([]model.Operator, error)
	FindByUserID(string) (*model.Operator, error)
	Create(*model.Operator) (*model.Operator, error)
	Update(*model.Operator) (*model.Operator, error)
	Delete(string) error
}

type OperatorImpl struct {
	*gorm.DB
}

func NewOperator(db *gorm.DB) Operator {
	return &OperatorImpl{db}
}

func (o *OperatorImpl) FindAll() ([]model.Operator, error) {
	operators := []model.Operator{}

	err := o.Find(&operators).Error

	return operators, err
}

func (o *OperatorImpl) FindByUserID(userID string) (*model.Operator, error) {
	operator := new(model.Operator)

	err := o.Where("user_id = ?", userID).Find(operator).Error

	return operator, err
}

func (o *OperatorImpl) Create(operator *model.Operator) (*model.Operator, error) {
	tx := o.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Create(operator).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return operator, tx.Commit().Error
}

func (o *OperatorImpl) Update(operator *model.Operator) (*model.Operator, error) {
	tx := o.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Where("id = ?", operator.ID).Save(operator).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return operator, tx.Commit().Error
}

func (o *OperatorImpl) Delete(id string) error {
	tx := o.Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("id = ?", id).Delete(&model.Operator{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
