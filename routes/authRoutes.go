package routes

import (
	"github.com/dev_mansoor/go-postgres-gorm/controllers"
	"github.com/dev_mansoor/go-postgres-gorm/middleware"
	"github.com/gin-gonic/gin"
)

func Authenticate(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/login", controllers.Login)
	incomingRoutes.POST("users/register", controllers.Register)
	incomingRoutes.GET("/users/validate", middleware.RequireAuth, controllers.Validate)
}
