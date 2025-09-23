package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("No .env file found, relying on environment variables: %v", err)
	}

	// Connect to the database
	if err := ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer DB.Close()

	router := gin.Default()

	// Define routes
	router.GET("/fish", getFish)
	router.GET("/fish/:id", getFishByID)
	router.POST("/fish", postFish)
	router.PUT("/fish/:id", updateFish)
	router.DELETE("/fish/:id", deleteFish)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run("localhost:" + port)
}

