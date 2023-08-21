package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"net/http"
	"sensor-metadata-api/internal/db"
	_ "sensor-metadata-api/internal/db"
	"strings"
	"time"
)

// CreateSensorMetadataHandler godoc
// @Summary      Create a new sensor metadata
// @Description  Create a new sensor metadata
// @Tags         create
// @Accept       json
// @Produce      json
// @Param        db_config.SensorMetadata   body     db.SensorMetadata   true    "SensorMetadata"
// @Success      201  {object}   interface{}
// @Failure      400  {object}  interface{}
// @Failure      500  {object}  interface{}
// @Router       /sensor-metadata [post]
func CreateSensorMetadataHandler(database db.SensorMetadataDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var sensor db.SensorMetadata
		if err := c.BodyParser(&sensor); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"code":    http.StatusBadRequest,
				"payload": map[string]string{"error": "invalid JSON"},
			})
		}

		if sensor.Name == "" || sensor.Location == (db.Location{}) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"code":    http.StatusBadRequest,
				"payload": map[string]string{"error": "sensor name and location are required"},
			})
		}

		sensor.CreatedAt = time.Now()
		sensor.UpdatedAt = time.Now()

		if err := database.CreateSensorMetadata(&sensor); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"payload": map[string]string{"error": "failed to insert sensor metadata: " + err.Error()},
			})
		}

		return c.Status(http.StatusCreated).JSON(fiber.Map{
			"code":    http.StatusCreated,
			"payload": map[string]string{"message": "successfully inserted sensor metadata"},
		})
	}
}

// GetSensorMetadataHandler godoc
// @Summary      Get info for a sensor
// @Description  Get info for a sensor
// @Tags         get
// @Accept       json
// @Produce      json
// @Param        name   path     string   true    "Sensor Name"
// @Success      200  {object}   db.SensorMetadata
// @Failure      404  {object}  interface{}
// @Failure      500  {object}  interface{}
// @Router       /sensor-metadata/{name} [get]
func GetSensorMetadataHandler(database db.SensorMetadataDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sensorName := c.Params("name")
		sn := strings.ToLower(sensorName)

		sensor, err := database.GetSensorMetadataByName(sn)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{
					"code":    http.StatusNotFound,
					"payload": map[string]string{"error": "sensor metadata not found"},
				})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"payload": map[string]string{"error": "failed to fetch sensor metadata"},
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"code":    http.StatusOK,
			"payload": sensor,
		})
	}
}

// UpdateSensorMetadataHandler godoc
// @Summary      Update sensor metadata
// @Description  Update sensor metadata
// @Tags         update
// @Accept       json
// @Produce      json
// @Param        name   path     string   true    "Sensor Name"
// @Param        name   body     db.SensorMetadata   true    "SensorMetadata"
// @Success      200  {object}   interface{}
// @Failure      404  {object}  interface{}
// @Failure      400  {object}  interface{}
// @Failure      500  {object}  interface{}
// @Router       /sensor-metadata/{name} [put]
func UpdateSensorMetadataHandler(database db.SensorMetadataDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sensorName := c.Params("name")
		sensorName = strings.ToLower(sensorName)

		sensor, err := database.GetSensorMetadataByName(sensorName)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{
					"code":    http.StatusNotFound,
					"payload": map[string]string{"error": "sensor metadata not found"},
				})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"payload": map[string]string{"error": "failed to fetch sensor metadata"},
			})
		}

		var updatedSensor db.SensorMetadata
		if err := c.BodyParser(&updatedSensor); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"code":    http.StatusBadRequest,
				"payload": map[string]string{"error": "invalid JSON"},
			})
		}

		if updatedSensor.Name != "" {
			sensor.Name = updatedSensor.Name
		}
		if updatedSensor.Location != (db.Location{}) {
			sensor.Location = updatedSensor.Location
		}
		if len(updatedSensor.Tags) > 0 {
			sensor.Tags = updatedSensor.Tags
		}

		sensor.UpdatedAt = time.Now()

		if err = database.CreateSensorMetadata(sensor); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"payload": map[string]string{"error": "failed to update sensor metadata"},
			})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"code":    http.StatusOK,
			"payload": map[string]string{"message": "successfully updated sensor metadata"},
		})
	}
}
