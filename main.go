package main

import (
	"fmt"

	"github.com/dev_mansoor/go-postgres-gorm/controllers"
	"github.com/dev_mansoor/go-postgres-gorm/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
}

func main() {
	fmt.Println("Hello, World! 123")

	router := gin.Default()
	router.POST("/posts", controllers.CreatePost)
	router.GET("/posts", controllers.PostIndex)
	router.GET("/posts/:id", controllers.PostShow)
	router.PUT("/posts/:id", controllers.PostUpdate)
	router.DELETE("/posts/:id", controllers.PostDelete)
	router.Run()
}
