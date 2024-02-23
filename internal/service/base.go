package service

import (
	"errors"
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gorm.io/gorm"
	"net/http"
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

type PagedResult[T any] struct {
	Data  []T   `json:"data"`
	Total int64 `json:"total"`
}

const PageSize = 10
const PageIndex = 1

// Pagination 分页查询
func Pagination[T any](db *gorm.DB, page int, size int, out []T) PagedResult[T] {
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
	db.Offset((page - 1) * size).Limit(size).Find(&out)

	return PagedResult[T]{
		Data:  out,
		Total: total,
	}
}

func ParamBadRequestResult(c *gin.Context, err error) {
	localize, _ := config.I18nLoc.LocalizeMessage(&i18n.Message{ID: "Error.param"})
	c.String(http.StatusBadRequest, localize)
	config.Log.Error(err.Error())
}

func BadRequestResult(c *gin.Context, messageId string, err error) {
	localize, _ := config.I18nLoc.LocalizeMessage(&i18n.Message{ID: messageId})
	c.String(http.StatusBadRequest, localize)
	config.Log.Error(err.Error())
}

func UnauthorizedResult(c *gin.Context, messageId string) {
	localize, _ := config.I18nLoc.LocalizeMessage(&i18n.Message{ID: messageId})
	c.String(http.StatusUnauthorized, localize)
}

func ConflictResult(c *gin.Context, messageId string) {
	localize, _ := config.I18nLoc.LocalizeMessage(&i18n.Message{ID: messageId})
	c.String(http.StatusConflict, localize)
}

func UpdateSuccessResult() domain.MessageWrapper {
	return SuccessResult("Success.update")
}

func DeleteSuccessResult() domain.MessageWrapper {
	return SuccessResult("Success.delete")
}

func SuccessResult(messageId string) domain.MessageWrapper {
	localize, _ := config.I18nLoc.LocalizeMessage(&i18n.Message{ID: messageId})
	return domain.NewMessageWrapper(localize)
}
