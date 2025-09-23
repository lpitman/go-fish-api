package main

import (
	"log"
	"os"
	"time"

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

	// Create repository and service instances
	fishRepo := NewFishRepository(DB)
	fishService := NewFishService(fishRepo)

	// Start the fish location simulation
	startSimulation(fishService)

	router := gin.Default()

	// Define routes
	router.GET("/fish", getFish(fishRepo))
	router.GET("/fish/:id", getFishByID(fishRepo))
	router.POST("/fish", postFish(fishRepo))
	router.PUT("/fish/:id", updateFish(fishRepo))
	router.DELETE("/fish/:id", deleteFish(fishRepo))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func startSimulation(fs *FishService) {
	log.Println("Creating timer for fish location updates every 10 seconds.")
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C{
			log.Println("--- Starting simulation cycle... ---")
			fs.UpdateFishLocations()
		}
	}()
}

