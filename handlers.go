package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var defaultLatitude = 44.692661
var defaultLongitude = -63.639532

// Responds with the list of all fish as JSON.
func getFish(repo *FishRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		fishes, err := repo.GetAll()
		if err != nil {
			log.Println("Error retrieving fish:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve fish data"})
			return
		}
		c.IndentedJSON(http.StatusOK, fishes)
	}
}

// getFishByID locates the fish whose ID value matches the id
func getFishByID(repo *FishRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		fish, err := repo.GetByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		if fish == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "fish not found"})
			return
		}
		c.IndentedJSON(http.StatusOK, fish)
	}
}

// postFish adds an fish from JSON received in the request body.
func postFish(repo *FishRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newFish Fish
		
		if err := c.BindJSON(&newFish); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}

		// Check if location data is missing and set to default if it is.
		if newFish.Location.Latitude == 0 && newFish.Location.Longitude == 0 {
			log.Println("Location data not provided. Setting to default Halifax location.")
			newFish.Location.Latitude = defaultLatitude
			newFish.Location.Longitude = defaultLongitude
		}

		newFish.ID = uuid.New().String()
		if err := repo.Create(&newFish); err != nil {
			log.Println("Error creating new fish:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create fish"})
			return
		}
		c.IndentedJSON(http.StatusCreated, newFish)
	}
}

// resetFishLocations back to the default Halifax location.
func resetFishLocations(repo *FishRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		fishes, err := repo.GetAll()
		if err != nil {
			log.Println("Error retrieving fish for reset:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve fish data"})
			return
		}

		for _, fish := range fishes {
			fish.Location.Latitude = defaultLatitude
			fish.Location.Longitude = defaultLongitude
			if _, err := repo.Update(fish); err != nil {
				log.Printf("Error resetting location for fish %s: %v\n", fish.ID, err)
			}
		}
		c.JSON(http.StatusOK, gin.H{"message": "All fish locations reset to default"})

	}
}
// updateFish updates a little fishy from JSON received in the request body.
func updateFish(repo *FishRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var updatedFish Fish
		if err := c.BindJSON(&updatedFish); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}
		updatedFish.ID = id // Ensure the ID from the URL is used

		rowsAffected, err := repo.Update(&updatedFish)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update fish"})
			return
		}
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "fish not found"})
			return
		}

		// Retrieve and return the full, updated object
		finalFish, _ := repo.GetByID(id)
		c.IndentedJSON(http.StatusOK, finalFish)
	}
}

// deleteFish removes a fish that has been eaten, by ID
func deleteFish(repo *FishRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		rowsAffected, err := repo.Delete(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete fish"})
			return
		}
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "fish not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

