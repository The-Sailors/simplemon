package data

import (
	"context"
	"database/sql"
	"log"
	"time"
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

type MonitorRepository interface {
	Create(ctx context.Context, monitor Monitor, log *log.Logger) error
}

func (m *MonitorModel) Create(ctx context.Context, monitor Monitor, log *log.Logger) (*Monitor, error) {
	var id int64
	err := m.DB.QueryRowContext(ctx, `
		INSERT INTO monitors (user_email, type, url, method, updated_at, body, headers, parameters, description, frequency_minutes, threshold_minutes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,  $11)
		RETURNING monitor_id`,
		monitor.UserEmail, monitor.MonitorType, monitor.URL, monitor.Method, monitor.UpdatedAt, monitor.Body, monitor.Headers, monitor.Parameters, monitor.Description, monitor.FrequencyMinutes, monitor.ThresholdMinutes).Scan(&id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	monitor.MonitorID = id
	return &monitor, nil
}
