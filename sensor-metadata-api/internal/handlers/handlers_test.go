package handlers

import (
	"bytes"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"sensor-metadata-api/internal/db"
	"strings"
	"testing"
	"time"
)

// Mock implementation for SensorMetadataDB
type MockSensorMetadataDB struct {
	mock.Mock
}

func (m *MockSensorMetadataDB) CreateSensorMetadata(sensor *db.SensorMetadata) error {
	args := m.Called(sensor)
	return args.Error(0)
}

func (m *MockSensorMetadataDB) GetSensorMetadataByName(name string) (*db.SensorMetadata, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*db.SensorMetadata), nil
}

func (m *MockSensorMetadataDB) UpdateSensorMetadata(sensor *db.SensorMetadata) error {
	args := m.Called(sensor)
	return args.Error(0)
}

func TestCreateSensorMetadataHandler_ValidInput(t *testing.T) {
	// Create mock database
	mockDB := new(MockSensorMetadataDB)

	// Set expectations
	mockDB.On("CreateSensorMetadata", mock.Anything).Return(nil)

	// Create handler instance
	handler := CreateSensorMetadataHandler(mockDB)

	// Create a new Fiber app
	app := fiber.New()

	// Define the route and use the handler
	app.Post("/sensor-metadata", handler)

	// Create a JSON request body with valid data
	requestBody := []byte(`{
		"name": "Sensor 1",
		"location": {
			"latitude": 40.123,
			"longitude": -75.456
		}
	}`)

	// Create a new POST request with the request body
	req := httptest.NewRequest(http.MethodPost, "/sensor-metadata", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}

	// Assertions
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Assert that the expectations were met
	mockDB.AssertExpectations(t)
}

func TestCreateSensorMetadataHandler_InvalidJSON(t *testing.T) {
	// Create mock database
	mockDB := new(MockSensorMetadataDB)

	// Create handler instance
	handler := CreateSensorMetadataHandler(mockDB)

	// Create a new Fiber app
	app := fiber.New()

	// Define the route and use the handler
	app.Post("/sensor-metadata", handler)

	// Create an invalid JSON request body
	requestBody := []byte(`{"invalid": "data"`)

	// Create a new POST request with the invalid request body
	req := httptest.NewRequest(http.MethodPost, "/sensor-metadata", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}

	// Assertions
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mockDB.AssertExpectations(t)
}

func TestCreateSensorMetadataHandler_MissingFields(t *testing.T) {
	// Create mock database
	mockDB := new(MockSensorMetadataDB)

	// Create handler instance
	handler := CreateSensorMetadataHandler(mockDB)

	// Create a new Fiber app
	app := fiber.New()

	// Define the route and use the handler
	app.Post("/sensor-metadata", handler)

	// Create a JSON request body with missing fields
	requestBody := []byte(`{
		"name": "Sensor 2"
	}`)

	// Create a new POST request with the request body
	req := httptest.NewRequest(http.MethodPost, "/sensor-metadata", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}

	// Assertions
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mockDB.AssertExpectations(t)
}

func TestCreateSensorMetadataHandler_DatabaseError(t *testing.T) {
	// Create mock database
	mockDB := new(MockSensorMetadataDB)

	// Set expectations for a database error
	mockDB.On("CreateSensorMetadata", mock.Anything).Return(errors.New("database error"))

	// Create handler instance
	handler := CreateSensorMetadataHandler(mockDB)

	// Create a new Fiber app
	app := fiber.New()

	// Define the route and use the handler
	app.Post("/sensor-metadata", handler)

	// Create a valid JSON request body
	requestBody := []byte(`{
		"name": "Sensor 3",
		"location": {
			"latitude": 40.123,
			"longitude": -75.456
		}
	}`)

	// Create a new POST request with the request body
	req := httptest.NewRequest(http.MethodPost, "/sensor-metadata", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	mockDB.AssertExpectations(t)
}

func TestGetSensorMetadataHandler_Success(t *testing.T) {
	// Create mock database
	mockDB := new(MockSensorMetadataDB)

	// Set expectations
	mockSensor := &db.SensorMetadata{Name: "Sensor1" /* ...other fields... */}
	mockDB.On("GetSensorMetadataByName", "sensor1").Return(mockSensor, nil)

	// Create handler instance with the mock database
	handler := GetSensorMetadataHandler(mockDB)

	// Create a new Fiber app
	app := fiber.New()

	// Define the route and use the handler
	app.Get("/sensor-metadata/:name", handler)

	// Create a new GET request for "/sensor-metadata/sensor1"
	req := httptest.NewRequest(http.MethodGet, "/sensor-metadata/sensor1", nil)

	// Create a response recorder
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}

	// Assertions
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Assert that the expectations were met
	mockDB.AssertExpectations(t)
}

func TestGetSensorMetadataHandler_NotFound(t *testing.T) {
	// Create mock database
	mockDB := new(MockSensorMetadataDB)

	// Set expectations for not found scenario
	mockDB.On("GetSensorMetadataByName", "unknownsensor").Return(nil, gorm.ErrRecordNotFound)

	// Create handler instance with the mock database
	handler := GetSensorMetadataHandler(mockDB)

	// Create a new Fiber app
	app := fiber.New()

	// Define the route and use the handler
	app.Get("/sensor-metadata/:name", handler)

	// Create a new GET request for "/sensor-metadata/unknownsensor"
	req := httptest.NewRequest(http.MethodGet, "/sensor-metadata/unknownsensor", nil)

	// Create a response recorder
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}

	// Assertions
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Assert that the expectations were met
	mockDB.AssertExpectations(t)
}

func TestGetSensorMetadataHandler_InternalServerError(t *testing.T) {
	// Create mock database
	mockDB := new(MockSensorMetadataDB)

	// Set expectations for internal server error scenario
	mockDB.On("GetSensorMetadataByName", "error").Return(nil, errors.New("error fetching data"))

	// Create handler instance with the mock database
	handler := GetSensorMetadataHandler(mockDB)

	// Create a new Fiber app
	app := fiber.New()

	// Define the route and use the handler
	app.Get("/sensor-metadata/:name", handler)

	// Create a new GET request for "/sensor-metadata/error"
	req := httptest.NewRequest(http.MethodGet, "/sensor-metadata/error", nil)

	// Create a response recorder
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockDB.AssertExpectations(t)
}

func TestUpdateSensorMetadataHandler(t *testing.T) {
	// Mock implementation of SensorMetadataDB
	mockDB := new(MockSensorMetadataDB)

	// Create the handler instance with the mock database
	handler := UpdateSensorMetadataHandler(mockDB)

	t.Run("Update_Success", func(t *testing.T) {
		// Mock sensor data
		sensorName := "sensor-1"
		mockSensor := &db.SensorMetadata{
			Name:      sensorName,
			Location:  db.Location{Latitude: 40.0, Longitude: -80.0},
			Tags:      []string{"tag1", "tag2"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Mock database method calls
		mockDB.On("GetSensorMetadataByName", sensorName).Return(mockSensor, nil)
		mockDB.On("CreateSensorMetadata", mock.Anything).Return(nil)

		// Create a new Fiber app
		app := fiber.New()

		// Define the route and use the handler
		app.Put("/sensor-metadata/:name", handler)

		// Create a new PUT request for updating sensor metadata
		payload := `{"name": "new-name", "location": {"latitude": 41.0, "longitude": -81.0}, "tags": ["new-tag"]}`
		req := httptest.NewRequest(http.MethodPut, "/sensor-metadata/"+sensorName, strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		resp, err := app.Test(req)
		assert.NoError(t, err)

		// Assertions
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Assert that the expectations were met
		mockDB.AssertExpectations(t)
	})

	t.Run("Invalid_JSON", func(t *testing.T) {
		// Mock sensor name
		sensorName := "sensor-1"

		// Mock database method call
		mockDB.On("GetSensorMetadataByName", sensorName).Return(nil, nil)

		// Create a new Fiber app
		app := fiber.New()

		// Define the route and use the handler
		app.Put("/sensor-metadata/:name", handler)

		// Create a new PUT request with invalid JSON payload
		req := httptest.NewRequest(http.MethodPut, "/sensor-metadata/"+sensorName, strings.NewReader("invalid-json"))
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		resp, err := app.Test(req)
		assert.NoError(t, err)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Assert that the expectations were met
		mockDB.AssertExpectations(t)
	})

	t.Run("Sensor_Not_Found", func(t *testing.T) {
		// Mock sensor name
		sensorName := "non-existent-sensor"

		// Mock database method call
		mockDB.On("GetSensorMetadataByName", sensorName).Return(nil, gorm.ErrRecordNotFound)

		// Create a new Fiber app
		app := fiber.New()

		// Define the route and use the handler
		app.Put("/sensor-metadata/:name", handler)

		// Create a new PUT request for updating non-existent sensor metadata
		req := httptest.NewRequest(http.MethodPut, "/sensor-metadata/"+sensorName, nil)

		// Create a response recorder
		resp, err := app.Test(req)
		assert.NoError(t, err)

		// Assertions
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		// Assert that the expectations were met
		mockDB.AssertExpectations(t)
	})
}
