package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/The-Sailors/simplemon/internal/data"
	"github.com/julienschmidt/httprouter"
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
		{
			name:   "Test createMonitorHandler monitor already exists",
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
				expectedStatusCode: 409,
				method:             "POST",
			},
			create: CreateReturn{
				monitor: nil,
				err:     data.ErrUniqueConstraintViolation,
			},
		},
		{
			name:   "Test createMonitorHandler not send the principal fields(URL, UserEmail, MonitorType, Method)",
			fields: Fields(initFields()),
			args: args{
				monitor: &data.Monitor{
					UpdatedAt:        time.Now(),
					Body:             "",
					Headers:          "",
					Parameters:       "",
					Description:      "",
					FrequencyMinutes: 1,
					ThresholdMinutes: 1,
				},
				expectedStatusCode: 400,
				method:             "POST",
			},
			create: CreateReturn{
				monitor: nil,
				err:     nil,
			},
		},
		{
			name:   "Test createMonitorHandler generic database error",
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
				expectedStatusCode: 500,
				method:             "POST",
			},
			create: CreateReturn{
				monitor: nil,
				err:     errors.New("generic database error"),
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

func TestApplication_getMonitorHandler(t *testing.T) {

	type args struct {
		expectedStatusCode int
		method             string
		monitor_id         string
	}
	type GetReturn struct {
		monitor *data.Monitor
		err     error
	}
	tests := []struct {
		name   string
		fields Fields
		args   args
		get    GetReturn // This is the return value of the Create method on the MonitorModelMock
	}{
		{
			name:   "Test getMonitorHandler",
			fields: Fields(initFields()),
			args: args{

				expectedStatusCode: 200,
				method:             "GET",
				monitor_id:         "1",
			},
			get: GetReturn{
				monitor: &data.Monitor{
					MonitorID:   1,
					URL:         "https://www.google.com",
					UserEmail:   "jojo@gmail.com",
					MonitorType: "jojo",
					Method:      "GET",
					UpdatedAt:   time.Now(),
					Body:        "",
					Headers:     "",
					Parameters:  "",
					Description: "",
				},
				err: nil,
			},
		},
		{
			name:   "Test getMonitorHandler monitor not found",
			fields: Fields(initFields()),
			args: args{
				expectedStatusCode: 404,
				method:             "GET",
				monitor_id:         "1",
			},
			get: GetReturn{
				monitor: nil,
				err:     data.ErrMonitorNotFound,
			},
		},
		{
			name:   "Test getMonitorHandler invalid monitor id",
			fields: Fields(initFields()),
			args: args{
				expectedStatusCode: 400,
				method:             "GET",
				monitor_id:         "invalid",
			},
			get: GetReturn{
				monitor: nil,
				err:     nil,
			},
		},
		{
			name:   "Test getMonitorHandler database generic error",
			fields: Fields(initFields()),
			args: args{
				expectedStatusCode: 500,
				method:             "GET",
				monitor_id:         "1",
			},
			get: GetReturn{
				monitor: nil,
				err:     errors.New("database generic error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testObj := data.NewMonitorModelMock()
			testObj.On("GetById", mock.Anything, mock.Anything, mock.Anything).Return(tt.get.monitor, tt.get.err)
			app := &Application{
				config: tt.fields.config,
				logger: tt.fields.logger,
				models: testObj,
			}
			//How to test query params: https://stackoverflow.com/questions/43502432/how-to-write-test-with-httprouter
			router := httprouter.New()
			router.HandlerFunc(tt.args.method, "/v1/monitors/:id", app.getMonitorHandler)

			req := httptest.NewRequest(tt.args.method, "/v1/monitors/"+tt.args.monitor_id, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			// handler := http.HandlerFunc(app.getMonitorHandler)
			// handler.ServeHTTP(w, req)
			if w.Code != tt.args.expectedStatusCode {
				t.Error(w.Body.String())
				t.Errorf("Expected status code %v, got %v", tt.args.expectedStatusCode, w.Code)
			}

		})
	}
}

func TestApplication_deleteMonitorHandler(t *testing.T) {

	type args struct {
		expectedStatusCode int
		method             string
		monitor_id         string
	}
	type GetReturn struct {
		monitor *data.Monitor
		err     error
	}
	type DeleteReturn struct {
		err error
	}
	tests := []struct {
		name   string
		fields Fields
		args   args
		get    GetReturn
		delete DeleteReturn // This is the return value of the Create method on the MonitorModelMock
	}{
		{
			name:   "Test deleteMonitorHandler deleted successfully",
			fields: Fields(initFields()),
			args: args{

				expectedStatusCode: 204,
				method:             "DELETE",
				monitor_id:         "1",
			},
			delete: DeleteReturn{
				err: nil,
			},
			get: GetReturn{
				monitor: nil,
				err:     nil,
			},
		},
		{
			name:   "Test deleteMonitorHandler monitor not found",
			fields: Fields(initFields()),
			args: args{
				expectedStatusCode: 404,
				method:             "DELETE",
				monitor_id:         "1",
			},
			delete: DeleteReturn{
				err: data.ErrMonitorNotFound,
			},
			get: GetReturn{
				monitor: nil,
				err:     data.ErrMonitorNotFound,
			},
		},
		{
			name:   "Test deleteMonitorHandler invalid monitor id",
			fields: Fields(initFields()),
			args: args{
				expectedStatusCode: 400,
				method:             "DELETE",
				monitor_id:         "invalid",
			},
			delete: DeleteReturn{
				err: nil,
			},
			get: GetReturn{
				monitor: nil,
				err:     nil,
			},
		},
		{
			name:   "Test deleteMonitorHandler Delete database generic error",
			fields: Fields(initFields()),
			args: args{
				expectedStatusCode: 500,
				method:             "DELETE",
				monitor_id:         "1",
			},
			delete: DeleteReturn{
				err: errors.New("database generic error"),
			},
			get: GetReturn{
				monitor: nil,
				err:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testObj := data.NewMonitorModelMock()
			//The delete handler calls the get function before deleting the verify if the monitor exists
			testObj.On("GetById", mock.Anything, mock.Anything, mock.Anything).Return(tt.get.monitor, tt.get.err)
			testObj.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(tt.delete.err)
			app := &Application{
				config: tt.fields.config,
				logger: tt.fields.logger,
				models: testObj,
			}
			//How to test query params: https://stackoverflow.com/questions/43502432/how-to-write-test-with-httprouter
			router := httprouter.New()
			router.HandlerFunc(tt.args.method, "/v1/monitors/:id", app.deleteMonitorHandler)

			req := httptest.NewRequest(tt.args.method, "/v1/monitors/"+tt.args.monitor_id, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.args.expectedStatusCode {
				t.Error(w.Body.String())
				t.Errorf("Expected status code %v, got %v", tt.args.expectedStatusCode, w.Code)
			}

		})
	}

}
