package main

import (
	"net/http"
	"os"

	"example.com/user/configs"
	"example.com/user/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()
	router.GET("/api/v1/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to User Microservice",
		})
	})
	//run database
	configs.ConnectDB()

	//routes
	routes.AuthRoutes(router)

	router.Run(":" + port)
}
