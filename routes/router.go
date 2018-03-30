package routes

import "github.com/gin-gonic/gin"
import (
	"../controller"
)

var router = gin.Default()
func InitRouter() {
	r := gin.Default()
	v1 := r.Group("api/v1/todos")
	{
		v1.GET("/", controller.FetchAllTodo1Handler)
		//v1.GET("/create", controller.CreateTodo)
		v1.GET("/:id", controller.FetchSingleTodo)
	}

	r.Run(":9898")
}