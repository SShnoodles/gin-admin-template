package router

import (
	"gin-admin-template/internal/api"
	"gin-admin-template/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetApiRouter(router *gin.Engine) {
	// login
	loginRouter := router.Group("login")
	loginRouter.Use(middleware.CORS(), middleware.Limit(1))
	{
		loginRouter.POST("account", api.Login)
		loginRouter.POST("captcha", api.Captcha)
	}
	// user
	usersRouter := router.Group("users")
	usersRouter.Use(middleware.CORS(), middleware.Auth())
	{
		usersRouter.GET("", api.GetUsers)
		usersRouter.GET(":id", api.GetUser)
		usersRouter.POST("", api.CreateUser)
		usersRouter.PUT(":id", api.UpdateUser)
		usersRouter.DELETE(":id", api.DeleteUser)
		usersRouter.GET(":id/roles", api.GetUserRoles)
		usersRouter.PUT(":id/enabled", api.EnabledUser)
	}
	// org
	orgsRouter := router.Group("orgs")
	orgsRouter.Use(middleware.CORS(), middleware.Auth())
	{
		orgsRouter.GET("", api.GetOrgs)
		orgsRouter.GET(":id", api.GetOrg)
		orgsRouter.POST("", api.CreateOrg)
		orgsRouter.PUT(":id", api.UpdateOrg)
		orgsRouter.DELETE(":id", api.DeleteOrg)
		orgsRouter.GET(":id/menus", api.GetOrgMenus)
	}
	// role
	roleRouter := router.Group("roles")
	roleRouter.Use(middleware.CORS(), middleware.Auth())
	{
		roleRouter.GET("", api.GetRoles)
		roleRouter.GET(":id", api.GetRole)
		roleRouter.POST("", api.CreateRole)
		roleRouter.PUT(":id", api.UpdateRole)
		roleRouter.DELETE(":id", api.DeleteRole)
		roleRouter.GET(":id/menus", api.GetRoleMenus)
		roleRouter.GET("orgs/:orgId", api.GetOrgRoles)
	}
	// menu
	menuRouter := router.Group("menus")
	menuRouter.Use(middleware.CORS(), middleware.Auth())
	{
		menuRouter.GET("", api.GetMenus)
		menuRouter.GET(":id", api.GetMenu)
		menuRouter.POST("", api.CreateMenu)
		menuRouter.PUT(":id", api.UpdateMenu)
		menuRouter.DELETE(":id", api.DeleteMenu)
		menuRouter.GET(":id/resources", api.GetMenuResources)
	}
	// resource
	resourceRouter := router.Group("resources")
	resourceRouter.Use(middleware.CORS(), middleware.Auth())
	{
		resourceRouter.GET("", api.GetResources)
	}
}
