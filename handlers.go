package main

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

type fish struct {
	ID     string  `json:"id"`
	Species string  `json:"species"`
	TrackingInfo  string     `json:"trackingInfo"`
	WeightKG  float64     `json:"weightKG"`
}

// fishData slice to seed our fish tracking data
// Later we'll put this in sqlite
var fishData = []fish{
	{ID: "1", Species: "Salmon", TrackingInfo: "Device-001", WeightKG: 4.5},
	{ID: "2", Species: "Tuna", TrackingInfo: "Device-002", WeightKG: 3.2},
	{ID: "3", Species: "Trout", TrackingInfo: "Device-003", WeightKG: 2.8},
}

// Responds with the list of all fish as JSON.
func getFish(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, fishData)
}

// getFishByID locates the fish whose ID value matches the id
func getFishByID(c *gin.Context) {
	id := c.Param("id")

	for _, f:= range fishData {
		if f.ID == id {
			c.IndentedJSON(http.StatusOK, f)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "fish not found"})
}

// postFish adds an fish from JSON received in the request body.
func postFish(c *gin.Context) {
	var newFish fish

	if err := c.BindJSON(&newFish); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a new ID for the fishy
	newID := strconv.Itoa(len(fishData) + 1)
	newFish.ID = newID
	fishData = append(fishData, newFish)
	c.IndentedJSON(http.StatusCreated, newFish)
}

// updateFish updates a little fishy from JSON received in the request body.
func updateFish(c *gin.Context) {
	id := c.Param("id")
	var updatedFish fish

	if err := c.BindJSON(&updatedFish); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// find fishy and update it
	for i, f:= range fishData {
		if f.ID == id {
			fishData[i].Species = updatedFish.Species
			fishData[i].TrackingInfo = updatedFish.TrackingInfo
			fishData[i].WeightKG = updatedFish.WeightKG
			c.IndentedJSON(http.StatusOK, fishData[i])
			return
		}
	}

	// No fishy found, return 404
	c.JSON(http.StatusNotFound, gin.H{"message": "fish not found"})
}

// deleteFish removes a fish that has been eaten, by ID
func deleteFish(c *gin.Context) {
	id := c.Param("id")

	for i, f:= range fishData {
		if f.ID == id {
			fishData = append(fishData[:i], fishData[i+1:]...)
			c.IndentedJSON(http.StatusNoContent, gin.H{"message": "bye bye fishy!"})
			return
		}
	}

	// No fishy found, return 404
	c.JSON(http.StatusNotFound, gin.H{"message": "fish not found"})
}
