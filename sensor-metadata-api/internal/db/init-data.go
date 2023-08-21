package db

import (
	"fmt"
	"time"
)

func initData(conn *SensorMetadataDBImpl) error {
	var sensors = []SensorMetadata{
		{
			Name:        "proximity",
			Description: "It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout.",
			Location: Location{
				Latitude:  40.25437,
				Longitude: -76.87133,
			},
			Tags:      []string{"tag1", "tag2"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:        "pressure",
			Description: "The point of using Lorem Ipsum is that it has a more-or-less normal distribution of letters, as opposed to using 'Content here, content here', making it look like readable English.",
			Location: Location{
				Latitude:  40.47202,
				Longitude: -80.01342,
			},
			Tags:      []string{"tag3", "tag4"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:        "capacitive",
			Description: "Many desktop publishing packages and web page editors now use Lorem Ipsum as their default model text, and a search for 'lorem ipsum' will uncover many web sites still in their infancy. Various versions have evolved over the years, sometimes by accident, sometimes on purpose (injected humour and the like)",
			Location: Location{
				Latitude:  39.9518,
				Longitude: -75.16845,
			},
			Tags:      []string{"tag6", "tag7"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	err := conn.db.Create(&sensors).Error
	if err != nil {
		return err
	}

	for _, sensor := range sensors {
		fmt.Println(sensor.ID)
	}

	return nil
}
