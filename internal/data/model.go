// tha data package is a layer that sits between the application and the database. It contains the
// models and the database migrations.
// Here in model file we will wrap the database models and expose them to the application using
// the NewModels function.
package data

import "database/sql"

type Models struct {
	Monitor *MonitorModel
}

type ModelsInterface interface {
	NewMonitorModel(db *sql.DB) *MonitorModel
}


func NewModels(db *sql.DB) Models {

	return Models{
		Monitor: NewMonitorModel(db),
	}
}

//TODO: Find a better place for this implementation
// Mock version of the models

type MockModels struct {
	Monitor *MonitorModelMock
}

func NewMockModels() MockModels {
	return MockModels{
		Monitor: NewMonitorModelMock(),
	}
}
