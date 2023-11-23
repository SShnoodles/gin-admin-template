package services

import (
	"errors"
	"gin-admin-template/initializations"
)

func FindById[T any](t *T, id int64) error {
	return initializations.DB.First(t, "id = ?", id).Error
}

func FindAll[T any](t *T) error {
	return initializations.DB.Find(t).Error
}

func Insert(i interface{}) error {
	return initializations.DB.Create(i).Error
}

func Update(i interface{}) error {
	return initializations.DB.Model(i).Updates(i).Error
}

func DeleteById(i interface{}, id int64) error {
	if id == 0 {
		return errors.New("id 为空！")
	}
	return initializations.DB.Delete(i, id).Error
}
