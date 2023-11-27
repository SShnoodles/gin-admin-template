package controllers

import (
	"gin-admin-template/initializations"
	"gin-admin-template/models"
	"gin-admin-template/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserQuery struct {
	services.PageInfo
	Username string `form:"username"`
	Mobile   string `form:"mobile"`
}

func GetUsers(c *gin.Context) {
	var q UserQuery
	err := c.ShouldBindQuery(&q)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	var users []models.User
	page := services.Pagination(initializations.DB, q.PageIndex, q.PageSize, &users)
	c.JSON(http.StatusOK, page)
}

func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var user models.User
	err = services.FindById(&user, int64(id))
	if err != nil {
		c.String(http.StatusBadRequest, "查询失败")
		return
	}
	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	user.Id = initializations.IdGenerate()
	err = services.Insert(user)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	c.JSON(http.StatusOK, models.NewIdWrapper(user.Id))
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	var user models.User
	err = c.ShouldBindJSON(&user)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	user.Id = int64(id)
	err = services.UpdateUser(user, false)
	if err != nil {
		c.String(http.StatusBadRequest, "更新失败")
		return
	}
	c.JSON(http.StatusOK, models.NewMessageWrapper("更新成功"))
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	err = services.DeleteById(models.User{}, int64(id))
	if err != nil {
		c.String(http.StatusBadRequest, "删除失败")
		return
	}
	c.JSON(http.StatusOK, models.NewMessageWrapper("删除成功"))
}
