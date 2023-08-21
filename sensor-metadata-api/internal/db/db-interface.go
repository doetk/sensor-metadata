package db

type SensorMetadataDB interface {
	CreateSensorMetadata(sensor *SensorMetadata) error
	GetSensorMetadataByName(name string) (*SensorMetadata, error)
	UpdateSensorMetadata(sensor *SensorMetadata) error
}
