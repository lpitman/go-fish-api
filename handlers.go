package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`	
}

type fish struct {
	ID     string  `json:"id"`
	Species string  `json:"species"`
	TrackingInfo  string     `json:"trackingInfo"`
	WeightKG  float64     `json:"weightKG"`
	Location  Location  `json:"location"`
}

// Responds with the list of all fish as JSON.
func getFish(c *gin.Context) {
	rows, err := DB.Query("SELECT id, species, tracking_info, weight_kg, latitude, longitude FROM fish")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var fishList []fish
	for rows.Next() {
		var f fish
		if err := rows.Scan(
			&f.ID, 
			&f.Species, 
			&f.TrackingInfo, 
			&f.WeightKG, 
			&f.Location.Latitude, 
			&f.Location.Longitude); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fishList = append(fishList, f)
	}
	c.IndentedJSON(http.StatusOK, fishList)
}

// getFishByID locates the fish whose ID value matches the id
func getFishByID(c *gin.Context) {
	id := c.Param("id")

	var f fish
	row := DB.QueryRow("SELECT id, species, tracking_info, weight_kg, latitude, longitude FROM fish WHERE id = ?", id) 
	if err := row.Scan(
		&f.ID, 
		&f.Species, 
		&f.TrackingInfo, 
		&f.WeightKG,
		&f.Location.Latitude,
		&f.Location.Longitude); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "fish not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, f)
}

// postFish adds an fish from JSON received in the request body.
func postFish(c *gin.Context) {
	var f fish
	if err := c.BindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := DB.Exec("INSERT INTO fish (species, tracking_info, weight_kg, latitude, longitude) VALUES (?, ?, ?, ?, ?)", 
		f.Species, f.TrackingInfo, f.WeightKG, f.Location.Latitude, f.Location.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, f)
}

// updateFish updates a little fishy from JSON received in the request body.
func updateFish(c *gin.Context) {
	id := c.Param("id")
	var f fish

	if err := c.BindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := DB.Exec("UPDATE fish SET species = ?, tracking_info = ?, weight_kg, latitude, longitude = ? WHERE id = ?", 
		f.Species, f.TrackingInfo, f.WeightKG, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "fish not found"})
		return
	}

	// Return the updated fish
	row := DB.QueryRow("SELECT id, species, tracking_info, weight_kg FROM fish WHERE id, latitude, longitude = ?", id)
	if err := row.Scan(
		&f.ID, 
		&f.Species, 
		&f.TrackingInfo, 
		&f.WeightKG,
		&f.Location.Latitude,
		&f.Location.Longitude); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.IndentedJSON(http.StatusOK, f)
}

// deleteFish removes a fish that has been eaten, by ID
func deleteFish(c *gin.Context) {
	id := c.Param("id")

	result, err := DB.Exec("DELETE FROM fish WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "fish not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
