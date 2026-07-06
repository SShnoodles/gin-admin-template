package web

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed dist
var distFS embed.FS

var distRoot = mustSub(distFS, "dist")

func RegisterRoutes(engine *gin.Engine) {
	engine.GET("/", serveIndex)
	engine.NoRoute(serveStaticOrIndex)
}

func mustSub(fsys fs.FS, dir string) fs.FS {
	subFS, err := fs.Sub(fsys, dir)
	if err != nil {
		panic(err)
	}
	return subFS
}

func serveIndex(c *gin.Context) {
	content, err := fs.ReadFile(distRoot, "index.html")
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", content)
}

func serveStaticOrIndex(c *gin.Context) {
	if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
		c.Status(http.StatusNotFound)
		return
	}
	if isAPIPath(c.Request.URL.Path) {
		c.Status(http.StatusNotFound)
		return
	}

	filePath := cleanRequestPath(c.Request.URL.Path)
	if fileExists(filePath) {
		c.FileFromFS("/"+filePath, http.FS(distRoot))
		return
	}

	serveIndex(c)
}

func cleanRequestPath(requestPath string) string {
	filePath := strings.TrimPrefix(requestPath, "/")
	if filePath == "" {
		return "index.html"
	}
	return path.Clean(filePath)
}

func fileExists(filePath string) bool {
	info, err := fs.Stat(distRoot, filePath)
	return err == nil && !info.IsDir()
}

func isAPIPath(requestPath string) bool {
	return requestPath == "/api" || strings.HasPrefix(requestPath, "/api/")
}
