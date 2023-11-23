package controllers

import (
	"github.com/gin-gonic/gin"
	"os"
)

func HtmlHandler(c *gin.Context) {
	dir, _ := os.Getwd()
	content, err := os.ReadFile(dir + "/web/dist/index.html")
	if err != nil {
		c.Writer.WriteString("not found")
		c.Writer.WriteHeader(404)
		return
	}
	c.Writer.Header().Add("Accept", "text/html")
	c.Writer.Header().Add("Content-Type", "text/html")
	c.Writer.Write(content)
	c.Writer.WriteHeader(200)
	c.Writer.Flush()
}
