package data

import "database/sql"

type Models struct {
	Monitor MonitorModel
}

func newModels(db *sql.DB) Models {
	return Models{
		Monitor: MonitorModel{DB: db},
	}
}
