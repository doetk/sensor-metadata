package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sensor-metadata-api/config"
)

func InitDb(cfg *config.DBConfig) (*SensorMetadataDBImpl, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
		cfg.SSLMode,
	)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Enable UUID Generator V4
	conn.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	// Auto-migrate the table
	_ = conn.Migrator().DropTable(&SensorMetadata{})
	err = conn.AutoMigrate(&SensorMetadata{})
	if err != nil {
		return nil, err
	}

	// Initialize the database instance
	db := &SensorMetadataDBImpl{db: conn}

	// Initialize initial data
	err = initData(db)
	if err != nil {
		return nil, err // Return error if initialization fails
	}

	return db, nil
}
