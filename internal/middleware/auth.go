package middleware

import (
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"gin-admin-template/internal/service"
	"gin-admin-template/internal/util"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verify jwt
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			service.UnauthorizedResult(c, "Error.unAuth")
			c.Abort()
			return
		}

		claims, err := util.ValidateUserToken(token, config.AppConfig.Jwt.Secret)
		if err != nil {
			service.UnauthorizedResult(c, "Error.unAuth")
			c.Abort()
			return
		}
		id, _ := claims.GetSubject()
		userId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			service.UnauthorizedResult(c, "Error.unAuth")
			c.Abort()
			return
		}
		c.Set("UserId", userId)

		// Verify user
		var user domain.User
		err = service.FindById(&user, userId)
		if err != nil {
			service.UnauthorizedResult(c, "Error.unAuth")
			c.Abort()
			return
		}
		if user == (domain.User{}) {
			service.UnauthorizedResult(c, "Error.unAuth")
			c.Abort()
			return
		}
		if !user.Enabled {
			service.UnauthorizedResult(c, "Error.unAuth")
			c.Abort()
			return
		}
		// ignore super admin
		if user.Username == "superadmin" {
			return
		}
		// Verify user resources
		if config.AppConfig.Verification.ResourceEnabled {
			resources, err := service.FindResourcesByUserId(userId)
			if err != nil {
				service.UnauthorizedResult(c, "Error.unAuth")
				c.Abort()
				return
			}

			if !isMatched(resources, c.Request.RequestURI, c.Request.Method) {
				service.UnauthorizedResult(c, "Error.unAuth")
				c.Abort()
			}
		}
	}
}

func isMatched(resources []domain.Resource, currentPath string, currentMethod string) bool {
	for _, r := range resources {
		if strings.EqualFold(r.Method, currentMethod) && isSamePattern(r.Path, currentPath) {
			return true
		}
	}
	return false
}

func isSamePattern(pattern, currentPath string) bool {
	paramsIndex := strings.Index(currentPath, "?")
	if paramsIndex != -1 {
		currentPath = currentPath[0:paramsIndex]
	}
	patternParts := strings.Split(pattern, "/")
	currentPathParts := strings.Split(currentPath, "/")

	if len(patternParts) != len(currentPathParts) {
		return false
	}

	for i, part := range patternParts {
		if part != currentPathParts[i] && !strings.HasPrefix(part, "{") && !strings.HasSuffix(part, "}") {
			return false
		}
	}

	return true
}
