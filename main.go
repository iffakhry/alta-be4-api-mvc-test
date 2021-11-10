package main

import (
	"alta/be4/mvc/config"
	"alta/be4/mvc/middlewares"
	"alta/be4/mvc/routes"
)

func main() {
	config.InitDB()
	e := routes.New()
	//implement LogMiddleware
	middlewares.LogMiddlewares(e)
	e.Logger.Fatal(e.Start(":8080"))
}
