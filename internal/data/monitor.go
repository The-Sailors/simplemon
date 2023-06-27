// This file contains the Monitor struct and the MonitorRepository interface, that are used to
// define the struct that will be managed and the functions to interact with this struct.
package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/rs/zerolog"
)

type Monitor struct {
	MonitorID        int64     `json:"monitor_id" `
	UserEmail        string    `json:"user_email"`
	MonitorType      string    `json:"type"`
	URL              string    `json:"url"`
	Method           string    `json:"method"`
	UpdatedAt        time.Time `json:"updated_at"`
	Body             string    `json:"body"`
	Headers          string    `json:"headers"`
	Parameters       string    `json:"parameters"`
	Description      string    `json:"description"`
	FrequencyMinutes int       `json:"frequency_minutes"`
	ThresholdMinutes int       `json:"threshold_minutes"`
}

type MonitorModel struct {
	DB *sql.DB
}

func NewMonitorModel(db *sql.DB) *MonitorModel {
	return &MonitorModel{DB: db}
}

type MonitorInterface interface {
	Create(ctx context.Context, monitor Monitor, log zerolog.Logger) (*Monitor, error)
	GetById(ctx context.Context, id int64, log zerolog.Logger) (*Monitor, error)
	Delete(ctx context.Context, id int64, log zerolog.Logger) error
	GetAll(ctx context.Context, log zerolog.Logger) ([]Monitor, error)
}

var (
	ErrUniqueConstraintViolation = errors.New("unique constraint violation")
	ErrMonitorNotFound           = errors.New("monitor not found")
)

func (m *MonitorModel) GetAll(ctx context.Context, log zerolog.Logger) ([]Monitor, error) {
	log.Info().Msg("Getting all monitors")
	rows, err := m.DB.QueryContext(ctx, `
		SELECT monitor_id, user_email, type, url, method, updated_at, body, headers, parameters, description, frequency_minutes, threshold_minutes
		FROM monitors`)
	if err != nil {
		log.Err(err).Msg("Error getting all monitors")
		return nil, err
	}
	defer rows.Close()

	monitors := []Monitor{}

	for rows.Next() {
		var monitor Monitor
		err := rows.Scan(&monitor.MonitorID, &monitor.UserEmail, &monitor.MonitorType, &monitor.URL, &monitor.Method, &monitor.UpdatedAt, &monitor.Body, &monitor.Headers, &monitor.Parameters, &monitor.Description, &monitor.FrequencyMinutes, &monitor.ThresholdMinutes)
		if err != nil {
			log.Err(err).Msg("Error scanning rows")
			return nil, err
		}
		monitors = append(monitors, monitor)
	}

	return monitors, nil
}

func (m *MonitorModel) Delete(ctx context.Context, id int64, log zerolog.Logger) error {
	log.Info().Msg("Deleting monitor")
	_, err := m.DB.ExecContext(ctx, `
		DELETE FROM monitors
		WHERE monitor_id = $1`,
		id)
	if err != nil {
		log.Err(err).Msg("Error deleting monitor")
		return err
	}
	return nil

}

func (m *MonitorModel) Create(ctx context.Context, monitor Monitor, log zerolog.Logger) (*Monitor, error) {
	log.Info().Msg("Creating monitor")
	var id int64
	var psqlErr *pq.Error

	err := m.DB.QueryRowContext(ctx, `
		INSERT INTO monitors (user_email, type, url, method, updated_at, body, headers, parameters, description, frequency_minutes, threshold_minutes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,  $11)
		RETURNING monitor_id`,
		monitor.UserEmail, monitor.MonitorType, monitor.URL, monitor.Method, monitor.UpdatedAt, monitor.Body, monitor.Headers, monitor.Parameters, monitor.Description, monitor.FrequencyMinutes, monitor.ThresholdMinutes).Scan(&id)
	if err != nil {
		log.Err(err).Msg("Error creating monitor")
		//if erro is pq: duplicate key value violates unique constraint "monitors_pkey"
		//then return a custom error
		if errors.As(err, &psqlErr) && psqlErr.Code == "23505" { // 23505 is unique_violation
			return nil, ErrUniqueConstraintViolation
		} else {

			return nil, err
		}
	}
	monitor.MonitorID = id
	return &monitor, nil
}

func (m *MonitorModel) GetById(ctx context.Context, id int64, log zerolog.Logger) (*Monitor, error) {
	log.Info().Msg("Getting monitor by id")
	var monitor Monitor
	err := m.DB.QueryRowContext(ctx, `
		SELECT monitor_id, user_email, type, url, method, updated_at, body, headers, parameters, description, frequency_minutes, threshold_minutes
		FROM monitors
		WHERE monitor_id = $1`,
		id).Scan(&monitor.MonitorID, &monitor.UserEmail, &monitor.MonitorType, &monitor.URL, &monitor.Method, &monitor.UpdatedAt, &monitor.Body, &monitor.Headers, &monitor.Parameters, &monitor.Description, &monitor.FrequencyMinutes, &monitor.ThresholdMinutes)
	if err != nil {
		//verify if the error is pq: no rows in result set
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMonitorNotFound
		} else {
			log.Err(err).Msg("Error getting monitor by id")
			return nil, err
		}

	}
	return &monitor, nil
}
