package db

import (
	"gorm.io/gorm"
)

type SensorMetadataDBImpl struct {
	db *gorm.DB
}

func NewSensorMetadataDB(db *gorm.DB) *SensorMetadataDBImpl {
	return &SensorMetadataDBImpl{db: db}
}

func (d *SensorMetadataDBImpl) CreateSensorMetadata(sensor *SensorMetadata) error {
	return d.db.Create(sensor).Error
}

func (d *SensorMetadataDBImpl) GetSensorMetadataByName(name string) (*SensorMetadata, error) {
	var sensor SensorMetadata
	if err := d.db.Where("name = ?", name).First(&sensor).Error; err != nil {
		return nil, err
	}

	return &sensor, nil
}

func (d *SensorMetadataDBImpl) UpdateSensorMetadata(sensor *SensorMetadata) error {
	return d.db.Save(sensor).Error
}
