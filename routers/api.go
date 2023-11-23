package routers

import (
	"gin-admin-template/controllers"
	"gin-admin-template/middlewares"
	"github.com/gin-gonic/gin"
)

func SetApiRouter(router *gin.Engine) {
	// login
	loginRouter := router.Group("/login")
	loginRouter.Use(middlewares.CORS(), middlewares.Limit(1))
	loginRouter.POST("/", controllers.Login)
	// user
	usersRouter := router.Group("/users")
	usersRouter.Use(middlewares.CORS())
	{
		usersRouter.GET("/", controllers.GetUsers)
		usersRouter.GET("/:id", controllers.GetUser)
		usersRouter.POST("/", controllers.CreateUser)
		usersRouter.PUT("/:id", controllers.UpdateUser)
		usersRouter.DELETE("/:id", controllers.DeleteUser)
	}
	// org
	orgsRouter := router.Group("/orgs")
	orgsRouter.Use(middlewares.CORS(), middlewares.Auth())
	{
		orgsRouter.GET("/", controllers.GetOrgs)
		orgsRouter.GET("/:id", controllers.GetOrg)
		orgsRouter.POST("/", controllers.CreateOrg)
		orgsRouter.PUT("/:id", controllers.UpdateOrg)
		orgsRouter.DELETE("/:id", controllers.DeleteOrg)
	}
	// role
	roleRouter := router.Group("/roles")
	roleRouter.Use(middlewares.CORS(), middlewares.Auth())
	{
		roleRouter.GET("/", controllers.GetRoles)
		roleRouter.GET("/:id", controllers.GetRole)
		roleRouter.POST("/", controllers.CreateRole)
		roleRouter.PUT("/:id", controllers.UpdateRole)
		roleRouter.DELETE("/:id", controllers.DeleteRole)
	}
	// menu
	menuRouter := router.Group("/menus")
	menuRouter.Use(middlewares.CORS(), middlewares.Auth())
	{
		menuRouter.GET("/", controllers.GetMenus)
		menuRouter.GET("/:id", controllers.GetMenu)
		menuRouter.POST("/", controllers.CreateMenu)
		menuRouter.PUT("/:id", controllers.UpdateMenu)
		menuRouter.DELETE("/:id", controllers.DeleteMenu)
	}
	// resource
	resourceRouter := router.Group("/resources")
	resourceRouter.Use(middlewares.CORS(), middlewares.Auth())
	{
		resourceRouter.GET("/", controllers.GetResources)
		resourceRouter.GET("/:id", controllers.GetResource)
		resourceRouter.POST("/", controllers.CreateResource)
		resourceRouter.PUT("/:id", controllers.UpdateResource)
		resourceRouter.DELETE("/:id", controllers.DeleteResource)
	}
}
