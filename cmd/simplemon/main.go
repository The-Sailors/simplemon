package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/The-Sailors/simplemon/internal/data"
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
	config Config      // All the configuration for the application
	logger *log.Logger // Generic logger for the application
	models data.Models // Models wraps all the application models.
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
			postgresURL:  getEnvWithDefault("POSTGRES_URL", "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable"),
			maxOpenConns: 5,
			maxIdleConns: 5,
		},
		env:  getEnvWithDefault("ENV", "development"),
		port: getEnvWithDefault("PORT", "8000"),
	}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := openDB(cfg, ctx)

	if err != nil {
		logger.Println("Cannot connect to database")
		logger.Fatal(err)
	}

	app := &Application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Printf("Starting server on %s", srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}
