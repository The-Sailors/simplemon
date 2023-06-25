package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/The-Sailors/simplemon/internal/data"
	"github.com/go-chi/httplog"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type Config struct {
	env      string
	port     string
	dbConfig struct {
		postgresURL  string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	logLevel  string
	logFormat string
}

func openDB(cfg Config, ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.dbConfig.postgresURL)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.dbConfig.maxOpenConns)

	db.SetMaxIdleConns(cfg.dbConfig.maxIdleConns)

	duration, err := time.ParseDuration(cfg.dbConfig.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

type Application struct {
	config Config                // All the configuration for the application
	logger zerolog.Logger        // Generic logger for the application
	models data.MonitorInterface // Models wraps all the application models.
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	cfg := Config{
		dbConfig: struct {
			postgresURL  string
			maxOpenConns int
			maxIdleConns int
			maxIdleTime  string
		}{
			postgresURL:  getEnvWithDefault("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
			maxOpenConns: 5,
			maxIdleConns: 5,
			maxIdleTime:  "15m",
		},
		env:       getEnvWithDefault("ENV", "development"),
		port:      getEnvWithDefault("PORT", "8000"),
		logLevel:  getEnvWithDefault("LOG_LEVEL", "info"),
		logFormat: getEnvWithDefault("LOG_FORMAT", "json"),
	}
	var json bool
	if cfg.logFormat == "json" {
		json = true
	} else {
		json = false
	}
	//structured logs
	logger := httplog.NewLogger("http", httplog.Options{
		JSON:     json,
		LogLevel: cfg.logLevel,
		Concise:  true,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := openDB(cfg, ctx)

	if err != nil {
		logger.Err(err).Msg("Cannot connect to database")
		logger.Fatal()
	}
	app := &Application{
		config: cfg,
		logger: logger,
		models: data.NewMonitorModel(db),
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Info().Msgf("Starting server on port %s", cfg.port)
	err = srv.ListenAndServe()
	logger.Fatal().Err(err)
}
