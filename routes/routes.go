package routes

import (
	"alta/be4/mvc/constants"
	"alta/be4/mvc/controllers"

	"github.com/labstack/echo/v4"
	echoMid "github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()
	e.GET("/users", controllers.GetUsersController)
	e.GET("/users/:id", controllers.GetOneUserController)
	e.POST("/users", controllers.CreateUserController)
	e.PUT("/users/:id", controllers.UpdateUserController)
	e.DELETE("/users/:id", controllers.DeleteUserController)
	e.POST("/login", controllers.LoginUsersController)

	//JWT Group
	eJWT := e.Group("/jwt")
	eJWT.Use(echoMid.JWT([]byte(constants.SECRET_JWT)))
	// eJWT.GET("/users/:id", controllers.GetOneUserController)
	eJWT.GET("/users/:id", controllers.GetUserDetailController)
	// eJWT.POST("/users", controllers.CreateUserController)

	//Basic Auth
	// eAuth := e.Group("")
	// eAuth.Use(echoMid.BasicAuth(mid.BasicAuthDB))
	// eAuth.DELETE("/users/:id", controllers.DeleteUserController)

	return e
}
