package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Define routes
	router.GET("/fish", getFish)
	router.GET("/fish/:id", getFishByID)
	router.POST("/fish", postFish)
	router.PUT("/fish/:id", updateFish)
	router.DELETE("/fish/:id", deleteFish)

	router.Run("localhost:8080")
}

