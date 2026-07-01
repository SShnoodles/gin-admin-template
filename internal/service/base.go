package service

import (
	"errors"
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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
		return errors.New("id is nil")
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

func Pagination[T any](db *gorm.DB, page int, size int, out []T) PagedResult[T] {
	if page == 0 {
		page = PageIndex
	}
	if size == 0 {
		size = PageSize
	}
	var total int64
	db.Model(out).Count(&total)

	db.Offset((page - 1) * size).Limit(size).Find(&out)

	return PagedResult[T]{
		Data:  out,
		Total: total,
	}
}

func ParamBadRequestResult(c *gin.Context) {
	localize, _ := config.I18nLoc.LocalizeMessage(&i18n.Message{ID: "Error.param"})
	Fail(c, http.StatusBadRequest, "400", localize)
}

func BadRequestResult(c *gin.Context, messageId string) {
	localize, _ := config.I18nLoc.LocalizeMessage(&i18n.Message{ID: messageId})
	Fail(c, http.StatusBadRequest, "400", localize)
}

func UnauthorizedResult(c *gin.Context, messageId string) {
	localize, _ := config.I18nLoc.LocalizeMessage(&i18n.Message{ID: messageId})
	Fail(c, http.StatusUnauthorized, "401", localize)
}

func ConflictResult(c *gin.Context, messageId string) {
	localize, _ := config.I18nLoc.LocalizeMessage(&i18n.Message{ID: messageId})
	Fail(c, http.StatusConflict, "409", localize)
}

func Ok[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, domain.SuccessResult(data))
}

func Fail(c *gin.Context, httpStatus int, errorCode string, message string) {
	c.JSON(httpStatus, domain.Result[any]{
		Success:      false,
		ErrorCode:    errorCode,
		ErrorMessage: message,
		ShowType:     2,
	})
}

func UpdateSuccessResult() domain.MessageWrapper {
	return SuccessMessage("Success.update")
}

func DeleteSuccessResult() domain.MessageWrapper {
	return SuccessMessage("Success.delete")
}

func SuccessMessage(messageId string) domain.MessageWrapper {
	localize, _ := config.I18nLoc.LocalizeMessage(&i18n.Message{ID: messageId})
	return domain.NewMessageWrapper(localize)
}
