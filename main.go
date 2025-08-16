package main

import (
	"fmt"

	"github.com/dev_mansoor/go-postgres-gorm/initializers"
	"github.com/dev_mansoor/go-postgres-gorm/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
}

func main() {
	fmt.Println("Hello, World! 123")
	// middleware.Text()

	router := gin.New()
	router.Use(gin.Logger())

	routes.PostRouter(router)
	routes.Authenticate(router)

	router.Run()
}
