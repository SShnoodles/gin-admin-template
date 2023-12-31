package service

import (
	"errors"
	"gin-admin-template/internal/config"
	"gorm.io/gorm"
)

func FindById[T any](t *T, id int64) error {
	return config.DB.First(t, "id = ?", id).Error
}

func FindByName[T any](t *T, name string) error {
	return config.DB.First(t, "name = ?", name).Error
}

func FindByCode[T any](t *T, code string) error {
	return config.DB.First(t, "code = ?", code).Error
}

func FindAll[T any](t *T) error {
	return config.DB.Find(t).Error
}

func Insert(i interface{}) error {
	return config.DB.Create(i).Error
}

func Update(i interface{}) error {
	return config.DB.Model(i).Updates(i).Error
}

func DeleteById(i interface{}, id int64) error {
	if id == 0 {
		return errors.New("id 为空！")
	}
	return config.DB.Delete(i, id).Error
}

type PageInfo struct {
	PageSize  int `form:"pageSize"`
	PageIndex int `form:"pageIndex"`
}

type PagedResult struct {
	Record any   `json:"record"`
	Count  int64 `json:"count"`
}

const PageSize = 10
const PageIndex = 1

// Pagination 分页查询
func Pagination(db *gorm.DB, page int, size int, out interface{}) PagedResult {
	if page == 0 {
		page = PageIndex
	}
	if size == 0 {
		size = PageSize
	}
	// 查询总数
	var total int64
	db.Model(out).Count(&total)

	// 分页查询当前页数据
	db.Offset((page - 1) * size).Limit(size).Find(out)

	return PagedResult{
		Record: out,
		Count:  total,
	}
}
