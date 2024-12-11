package bootstrap

import (
	"github.com/kholidss/movie-fest-skilltest/pkg/config"
	"github.com/kholidss/movie-fest-skilltest/pkg/logger"
	"github.com/kholidss/movie-fest-skilltest/pkg/util"
)

func RegistryLogger(cfg *config.Config) {
	loggerConfig := logger.Config{
		Environment: util.EnvironmentTransform(cfg.AppEnv),
		Debug:       cfg.AppDebug,
		Level:       cfg.LogLevel,
		ServiceName: cfg.AppName,
	}

	logger.Setup(loggerConfig)

	switch cfg.LogDriver {
	case logger.LogDriverLoki:
		logger.AddHook(cfg.LoggerConfig.WithLokiHook(cfg))
		break
	case logger.LogDriverGraylog:
		logger.AddHook(cfg.LoggerConfig.WithGraylogHook(cfg))
		break
	default:
		break
	}
}
