package routes

import (
	"github.com/dev_mansoor/go-postgres-gorm/controllers"
	"github.com/dev_mansoor/go-postgres-gorm/middleware"
	"github.com/gin-gonic/gin"
)

func PostRouter(router *gin.Engine) {
	router.Use(middleware.Authenticate)
	router.POST("/posts", middleware.RequireAuth, controllers.CreatePost)
	router.GET("/posts", middleware.RequireAuth, controllers.PostIndex)
	router.GET("/posts/:id", middleware.RequireAuth, controllers.PostShow)
	router.PUT("/posts/:id", middleware.RequireAuth, controllers.PostUpdate)
	router.DELETE("/posts/:id", middleware.RequireAuth, controllers.PostDelete)
}
