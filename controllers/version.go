package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const Version = "1.0.0"

func GetVersion(c *gin.Context) {
	c.String(http.StatusOK, Version)
}
