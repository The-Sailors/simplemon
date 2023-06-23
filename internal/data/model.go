// tha data package is a layer that sits between the application and the database. It contains the
// models and the database migrations.
// Here in model file we will wrap the database models and expose them to the application using
// the NewModels function.
package data

import "database/sql"

type Models struct {
	Monitor MonitorModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Monitor: MonitorModel{DB: db},
	}
}
