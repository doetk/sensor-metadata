package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"net/http"
	"sensor-metadata-api/internal/db"
	"sensor-metadata-api/internal/handlers"
	"sensor-metadata-api/internal/logger"
	"sensor-metadata-api/internal/version"
)

func (s *Server) SetupRoutes(database db.SensorMetadataDB) {

	s.app.Use(cors.New())

	// Serve Swagger Docs
	s.app.Get("/swagger/*", swagger.HandlerDefault)

	// main API v1 group
	api := s.app.Group(
		"/api/v1",
		logger.RequestId(),
		logger.WrapLogger(),
	)

	// monitor endpoint - /api/v1/monitor
	api.Get("/monitor", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(map[string]any{
			"code": http.StatusOK,
			"payload": map[string]string{
				"status":  "sensor-metadata-api is running",
				"version": version.ProductVersion,
			},
		})
	})

	// API V1 Group
	v1 := api.Group(
		"/sensor-metadata",
	)

	v1.Post("", handlers.CreateSensorMetadataHandler(database))
	v1.Get("/:name", handlers.GetSensorMetadataHandler(database))
	v1.Put("/:name", handlers.UpdateSensorMetadataHandler(database))
}
