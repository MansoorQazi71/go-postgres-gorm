package main

import (
	"github.com/dev_mansoor/go-postgres-gorm/initializers"
	"github.com/dev_mansoor/go-postgres-gorm/initializers/models"
)

func init() {
	// Load environment variables and connect to the database
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
}

func main() {
	initializers.DB.AutoMigrate(&models.Post{})
	initializers.DB.AutoMigrate(&models.User{})
}
