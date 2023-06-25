package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/The-Sailors/simplemon/internal/data"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
)

type Fields struct {
	config Config
	logger zerolog.Logger
}

func initFields() Fields {

	cfg := Config{
		env:       "dev",
		logLevel:  "error",
		logFormat: "text",
	}
	return Fields{
		config: cfg,
		logger: setupLog(cfg),
	}
}

func TestApplication_createMonitorHandler(t *testing.T) {

	type args struct {
		monitor            *data.Monitor
		expectedStatusCode int
		method             string
	}
	type CreateReturn struct {
		monitor *data.Monitor
		err     error
	}
	tests := []struct {
		name   string
		fields Fields
		args   args
		create CreateReturn // This is the return value of the Create method on the MonitorModelMock
	}{
		{
			name:   "Test createMonitorHandler",
			fields: Fields(initFields()),
			args: args{
				monitor: &data.Monitor{
					URL:              "https://www.google.com",
					UserEmail:        "jojo@gmail.com",
					MonitorType:      "jojo",
					Method:           "GET",
					UpdatedAt:        time.Now(),
					Body:             "",
					Headers:          "",
					Parameters:       "",
					Description:      "",
					FrequencyMinutes: 1,
					ThresholdMinutes: 1,
				},
				expectedStatusCode: 201,
				method:             "GET",
			},
			create: CreateReturn{
				monitor: &data.Monitor{
					MonitorID:        1,
					URL:              "https://www.google.com",
					UserEmail:        "jojo@gmail.com",
					MonitorType:      "jojo",
					Method:           "GET",
					UpdatedAt:        time.Now(),
					Body:             "",
					Headers:          "",
					Parameters:       "",
					Description:      "",
					FrequencyMinutes: 1,
					ThresholdMinutes: 1,
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testObj := data.NewMonitorModelMock()
			testObj.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(tt.create.monitor, tt.create.err)
			app := &Application{
				config: tt.fields.config,
				logger: tt.fields.logger,
				models: testObj,
			}
			monitorJson, err := json.Marshal(tt.args.monitor)
			if err != nil {
				t.Errorf("Error marshalling monitor: %v", err)
			}
			monitorString := string(monitorJson)
			req := httptest.NewRequest(tt.args.method, "/v1/monitors", strings.NewReader(monitorString))
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(app.createMonitorHandler)
			handler.ServeHTTP(w, req)
			if w.Code != tt.args.expectedStatusCode {
				t.Errorf("Expected status code %v, got %v", tt.args.expectedStatusCode, w.Code)
			}

		})
	}
}
