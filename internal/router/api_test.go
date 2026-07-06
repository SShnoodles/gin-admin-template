package router

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetApiRouterUsesAPIPrefix(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	SetApiRouter(engine)
	SetOtherRouter(engine)

	routes := map[string]bool{}
	for _, route := range engine.Routes() {
		routes[route.Method+" "+route.Path] = true
	}

	for _, route := range []string{
		"POST /api/login/account",
		"GET /api/users",
		"GET /api/menus",
		"GET /api/resources",
		"GET /api/project/version",
	} {
		if !routes[route] {
			t.Fatalf("expected route %s to be registered", route)
		}
	}
}
