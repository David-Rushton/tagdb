package tagdb

import (
	"time"

	"dev.azure.com/trayport/Hackathon/_git/Q/internal/logger"
)

const (
	defaultRollWalAfterBytes        = 10 * 1024 * 1024 // 10 MiB.
	defaultBackgroundTaskIntervalMs = 1_000            // 1 second.
)

// Configures the database.
type dbConfig struct {
	// The current WAL should roll at the opportunity after exceeding this size.
	rollWalAfterBytes int64

	// Background tasks are run at this interval.
	backgroundTaskInterval time.Duration
}

type dbConfigurer func(dbConfig *dbConfig) *dbConfig

// Applies default values to all configuration options.
func WithDefaultConfig() dbConfigurer {
	return func(dbConfig *dbConfig) *dbConfig {
		// Validation.
		if dbConfig == nil {
			logger.Panic("cannot configure database")
		}

		interval := time.Millisecond * defaultBackgroundTaskIntervalMs
		dbConfig.backgroundTaskInterval = interval
		dbConfig.rollWalAfterBytes = defaultRollWalAfterBytes

		return dbConfig
	}
}

// Defines the point at which the WAL will consider rolling.
func WithRollAfterBytes(value int64) dbConfigurer {
	return func(dbConfig *dbConfig) *dbConfig {
		// Validation.
		if dbConfig == nil {
			logger.Panic("cannot configure database")
		}

		if value <= 0 {
			logger.Panic("cannot configure database, rollAfterBytes must be great than 0")
		}

		dbConfig.rollWalAfterBytes = value

		return dbConfig
	}
}

// Defines the number of ticks between house keeping activities.
func WithBackgroundTaskIntervalMs(value int) dbConfigurer {
	return func(dbConfig *dbConfig) *dbConfig {
		// Validation.
		if dbConfig == nil {
			logger.Panic("cannot configure database")
		}

		if value <= 0 {
			logger.Panic("cannot configure database, houseKeepingTickIntervalMs must be great than 0")
		}

		if value < 100 {
			logger.Warnf("background refresh intervals of %d is below the minimum recommended value of 100", value)
		}

		interval := time.Millisecond * time.Duration(value)
		dbConfig.backgroundTaskInterval = interval

		return dbConfig
	}
}
