package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type FishService struct {
	repo *FishRepository
}

func NewFishService(repo *FishRepository) *FishService {
	return &FishService{repo: repo}
}

// UpdateFishLocations is the main simulation function.
func (s *FishService) UpdateFishLocations() {
	// Seed the random number generator.
	rand.Seed(time.Now().UnixNano())

	fishes, err := s.repo.GetAll()
	if err != nil {
		log.Println("Error fetching fish for update:", err)
		return
	}

	eatenFish := make(map[string]bool)

	// Collision and interaction logic
	for i := 0; i < len(fishes); i++ {
		if eatenFish[fishes[i].ID] {
			continue
		}
		for j := i + 1; j < len(fishes); j++ {
			if eatenFish[fishes[j].ID] {
				continue
			}

			if haversine(fishes[i].Location.Latitude, fishes[i].Location.Longitude, fishes[j].Location.Latitude, fishes[j].Location.Longitude) < 0.05 {
				log.Printf("Collision detected between Fish %s and Fish %s.\n", fishes[i].ID, fishes[j].ID)
				if fishes[i].Species == fishes[j].Species {
					// Mating
					newFish := Fish{
						ID:           uuid.New().String(),
						Species:      fishes[i].Species,
						TrackingInfo: "Offspring",
						WeightKG:     (fishes[i].WeightKG + fishes[j].WeightKG) / 2,
						Location: Location{
							Latitude:  (fishes[i].Location.Latitude + fishes[j].Location.Latitude) / 2,
							Longitude: (fishes[i].Location.Longitude + fishes[j].Location.Longitude) / 2,
						},
					}
					if err := s.repo.Create(&newFish); err != nil {
						log.Println("Error creating new fish after mating:", err)
					} else {
						log.Printf("Mating successful! New fish created with ID: %s\n", newFish.ID)
					}
				} else {
					// Eating
					var winner, loser *Fish
					if fishes[i].WeightKG > fishes[j].WeightKG {
						winner, loser = fishes[i], fishes[j]
					} else {
						winner, loser = fishes[j], fishes[i]
					}
					winner.WeightKG += loser.WeightKG
					if _, err := s.repo.Delete(loser.ID); err != nil {
						log.Println("Error deleting eaten fish:", err)
					} else {
						eatenFish[loser.ID] = true
						log.Printf("Fish %s was eaten by Fish %s.\n", loser.ID, winner.ID)
					}
				}
			}
		}
	}

	// Update location and weight for all remaining fish.
	for _, f := range fishes {
		if eatenFish[f.ID] {
			continue
		}
		f.Location.Latitude += (rand.Float64() - 0.5) * 0.01
		f.Location.Longitude += (rand.Float64() - 0.5) * 0.01
		f.WeightKG += 0.1
		if _, err := s.repo.Update(f); err != nil {
			log.Printf("Error updating fish %s: %v\n", f.ID, err)
		} else {
			log.Printf("Updated fish %s - New Loc: (%.4f, %.4f), New Wt: %.2f kg\n", f.ID, f.Location.Latitude, f.Location.Longitude, f.WeightKG)
		}
	}
	log.Println("--- Simulation cycle finished. ---")
}

// haversine calculates the distance between two lat/lon points in kilometers.
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in kilometers.
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*math.Pi/180.0)*math.Cos(lat2*math.Pi/180.0)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

