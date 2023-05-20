package routes

import (
	"example.com/user/controllers"
	"example.com/user/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.CORSMiddleware())
	incomingRoutes.POST("/users/signup", controllers.Signup())
	incomingRoutes.POST("/users/login", controllers.Login())
}
