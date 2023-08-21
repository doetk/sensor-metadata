package db

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

// SensorMetadata Sensor represents the sensor metadata structure
type SensorMetadata struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name        string         `gorm:"type:varchar(255); not null; unique"  json:"name"`
	Description string         `gorm:"type:varchar; not null; unique"  json:"description"`
	Location    Location       `gorm:"embedded" json:"location"`
	Tags        pq.StringArray `gorm:"type:text[]" json:"tags"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// Location represents the GPS position
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
