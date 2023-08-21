package main

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"sensor-metadata-api/config"
	_ "sensor-metadata-api/docs"
	db_config "sensor-metadata-api/internal/db"
	"sensor-metadata-api/internal/server"
	"syscall"
	"time"
)

// @title          Sensor Metadata API Application
// @version 2.0
// @description    This is a single-binary sensor-metadata-application
// @termsOfService  http://swagger.io/terms/

// @contact.name   Theo Doe
// @contact.url    https://www.linkedin.com/in/tkdoe/
// @contact.email  info.tkdoe@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath       /api/v1/
func main() {
	cfg := config.InitConfig(".")
	logger := setupLogging()
	defer logger.Sync()

	logger.Info("setting up db connection")
	db, err := db_config.InitDb(cfg.DBConfig)
	if err != nil {
		logger.Fatal("error setting up db connection: " + err.Error())
	}

	s := server.NewServer(fiber.Config{
		ReadBufferSize:        1 << 20,
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          10 * time.Second,
		IdleTimeout:           10 * time.Second,
		DisableStartupMessage: true,
	})

	s.SetupRoutes(db)

	go func() {
		logger.Info("server listener starting " + cfg.ServerConfig.Addr)
		err := s.Listen(cfg.ServerConfig.Addr)
		if err != nil {
			logger.Fatal(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	sig := <-c
	logger.Info("ending process " + sig.String() + " signal received")

	err = s.Shutdown()
	if err != nil {
		logger.Error("error shutting down server: " + err.Error())
	}
}

func setupLogging() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
	}

	return zap.Must(cfg.Build())
}
