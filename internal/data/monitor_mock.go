// The functions in this file are used to mock the database connection and the database model.
package data

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
)

type MonitorModelMock struct {
	mock.Mock
}

func NewMonitorModelMock() *MonitorModelMock {
	return &MonitorModelMock{}
}

type MockMonitorInterface interface {
	Create(ctx context.Context, monitor Monitor, log zerolog.Logger) (*Monitor, error)
}

func (m *MonitorModelMock) Create(ctx context.Context, monitor Monitor, log zerolog.Logger) (*Monitor, error) {
	args := m.Called(ctx, monitor, log)
	return args.Get(0).(*Monitor), args.Error(1)
}
