package main

import (
	"database/sql"
	"errors"
	"log"
)

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`	
}

type Fish struct {
	ID     string  `json:"id"`
	Species string  `json:"species"`
	TrackingInfo  string     `json:"trackingInfo"`
	WeightKG  float64     `json:"weightKG"`
	Location  Location  `json:"location"`
}

// FishRepository will handle db operations
type FishRepository struct {
	db *sql.DB
}

// Create a new FishRepository with a db connection
func NewFishRepository(db *sql.DB) *FishRepository {
	return &FishRepository{db: db}
}

func (r *FishRepository) GetAll() ([]*Fish, error) {
	rows, err := r.db.Query(`
		SELECT
			id,
			species,
			tracking_info,
			weight_kg,
			latitude,
			longitude
		FROM fish
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fishes []*Fish
	for rows.Next() {
		var f Fish
		if err := rows.Scan(
			&f.ID,
			&f.Species,
			&f.TrackingInfo,
			&f.WeightKG,
			&f.Location.Latitude,
			&f.Location.Longitude,
		); err != nil {
			return nil, err
		}
		fishes = append(fishes, &f)
	}
	return fishes, nil
}

// GetByID retrieves a single fish by its ID.
func (r *FishRepository) GetByID(id string) (*Fish, error) {
	var f Fish
	row := r.db.QueryRow(`
		SELECT 
			id, 
			species, 
			tracking_info, 
			weight_kg, 
			latitude, 
			longitude 
		FROM fish 
		WHERE id = ?
		`, id)
	err := row.Scan(
		&f.ID, 
		&f.Species, 
		&f.TrackingInfo, 
		&f.WeightKG, 
		&f.Location.Latitude, 
		&f.Location.Longitude)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil 
		}
		return nil, err 
	}
	return &f, nil
}

// Create inserts a new fish record into the database.
func (r *FishRepository) Create(f *Fish) error {
	log.Println("Creating fish with ID:", f.ID)
	log.Printf("Fish details: Species=%s, TrackingInfo=%s, WeightKG=%.2f, Location=(Lat: %.6f, Long: %.6f)\n", 
		f.Species, f.TrackingInfo, f.WeightKG, f.Location.Latitude, f.Location.Longitude)
	_, err := r.db.Exec(`
		INSERT INTO fish (
			id, 
			species, 
			tracking_info, 
			weight_kg, 
			latitude, 
			longitude
		) 
		VALUES (?, ?, ?, ?, ?, ?)
	`, f.ID, f.Species, f.TrackingInfo, f.WeightKG, f.Location.Latitude, f.Location.Longitude)
	
	return err
}

// Update modifies a fish's data in the database.
func (r *FishRepository) Update(f *Fish) (int64, error) {
	result, err := r.db.Exec(`
		UPDATE fish 
		SET 
			species = ?, 
			tracking_info = ?, 
			weight_kg = ?, 
			latitude = ?, 
			longitude = ? 
		WHERE id = ?
		`, f.Species, f.TrackingInfo, f.WeightKG, f.Location.Latitude, f.Location.Longitude, f.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete removes a fish from the database by its ID.
func (r *FishRepository) Delete(id string) (int64, error) {
	result, err := r.db.Exec(`
		DELETE 
		FROM fish 
		WHERE id = ?
		`, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

