package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/The-Sailors/simplemon/internal/data"
	"github.com/stretchr/testify/mock"
)

type Fields struct {
	config Config
	logger *log.Logger
}

func initFields() Fields {
	postgresURL := "postgres://postgres:postgres@localhost:5432/simplemon?sslmode=disable"
	maxOpenConns := 5
	maxIdleConns := 5

	return Fields{
		config: Config{
			env:  "dev",
			port: "8080",
			dbConfig: struct {
				postgresURL  string
				maxOpenConns int
				maxIdleConns int
				maxIdleTime  string
			}{
				postgresURL:  postgresURL,
				maxOpenConns: maxOpenConns,
				maxIdleConns: maxIdleConns,
			},
		},
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}

func TestApplication_createMonitorHandler(t *testing.T) {

	type fields struct {
		config Config
		logger *log.Logger
	}
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
		fields fields
		args   args
		create CreateReturn // This is the return value of the Create method on the MonitorModelMock
	}{
		{
			name:   "Test createMonitorHandler",
			fields: fields(initFields()),
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
