package configs

import (
	"time"

	"github.com/labstack/echo/v4/middleware"
)

func GetEchoLoggerConfig() middleware.LoggerConfig {
	logger := middleware.DefaultLoggerConfig

	logger.CustomTimeFormat = time.RFC3339
	logger.Format = "${time_rfc3339} (${status}) [${method}] ${uri} - ${latency_human}\n"

	return logger
}

func GetEchoRecoverConfig() middleware.RecoverConfig {
	recoverConfig := middleware.DefaultRecoverConfig

	// recoverConfig.DisablePrintStack = true

	return recoverConfig
}
