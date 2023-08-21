package logger

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

const (
	RequestIdCtxKey = "requestid"
	loggerCtxKey    = "logger"
)

func RequestId() fiber.Handler {
	return requestid.New()
}

func createLogger(c *fiber.Ctx) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: true,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		InitialFields: map[string]interface{}{
			"transaction_id": c.Locals(RequestIdCtxKey),
			"remote_addr":    c.IP(),
			"req_method":     c.Method(),
			"req_url":        string(c.Request().URI().RequestURI()),
		},
	}

	return zap.Must(config.Build())
}

func WrapLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqTime := time.Now()

		l := createLogger(c)
		defer l.Sync()

		c.Locals(loggerCtxKey, l)

		// call the next handler in the chain
		err := c.Next()
		if err != nil {
			l.Info(err.Error(),
				zap.Int64("response_time", time.Since(reqTime).Milliseconds()),
				zap.Int("response_size", len(c.Response().Body())),
				zap.Int("status_code", c.Response().StatusCode()),
			)
			return err
		}

		// log the response
		switch {
		// log the error response sent to the client
		case c.Response().StatusCode() != http.StatusOK && c.Response().StatusCode() != http.StatusFound:
			r := map[string]any{}
			if err = json.Unmarshal(c.Response().Body(), &r); err != nil {

				l.Error("error unmarshalling payload response: "+err.Error(),
					zap.Int64("response_time", time.Since(reqTime).Milliseconds()),
					zap.Int("response_size", len(c.Response().Body())),
					zap.Int("status_code", c.Response().StatusCode()),
				)
				return err
			}

			l.Error(fmt.Sprintf("%s", r["payload"]),
				zap.Int64("response_time", time.Since(reqTime).Milliseconds()),
				zap.Int("response_size", len(c.Response().Body())),
				zap.Int("status_code", c.Response().StatusCode()),
			)

		default: // log that we sent a response
			logger := l.Level()
			switch logger.Get() {
			case "Debug":
				l.Info(string(c.Response().Body()),
					zap.Int64("response_time", time.Since(reqTime).Milliseconds()),
					zap.Int("response_size", len(c.Response().Body())),
					zap.Int("status_code", c.Response().StatusCode()),
				)
			default:
				l.Info("sent response",
					zap.Int64("response_time", time.Since(reqTime).Milliseconds()),
					zap.Int("response_size", len(c.Response().Body())),
					zap.Int("status_code", c.Response().StatusCode()),
				)
			}
		}

		return nil
	}
}

func SetLoggerForRequest(c *fiber.Ctx, l *zap.Logger) {
	c.Locals(loggerCtxKey, l)
}
