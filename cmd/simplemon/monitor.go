package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/The-Sailors/simplemon/internal/data"
	"github.com/go-chi/httplog"
)

func (app *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info().Msg("Starting Healthcheck Handler")
	fmt.Fprintln(w, "Jojo: is awesome!")
	fmt.Fprintln(w, "environment:", app.config.env)
}

func (app *Application) createMonitorHandler(w http.ResponseWriter, r *http.Request) {
	log := httplog.LogEntry(r.Context())
	var monitor data.Monitor

	err := json.NewDecoder(r.Body).Decode(&monitor)
	if err != nil {
		log.Err(err).Msg("Error decoding the request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//verify if the primary keys fields user email, type, url and method are not empty
	if monitor.UserEmail == "" || monitor.MonitorType == "" || monitor.URL == "" || monitor.Method == "" {
		log.Err(nil).Msg("User email, type, url and method are required")
		http.Error(w, "User email, type, url and method are required", http.StatusBadRequest)
		return
	}
	//Create the monitor in the database
	createdMonitor, err := app.models.Create(r.Context(), monitor, log)
	if err != nil {
		log.Err(err).Msg("Error creating the monitor")
		//verify if the error is a unique constraint violation

		http.Error(w, "Database err", http.StatusInternalServerError)
		return
	}
	createdMonitorJson, err := json.Marshal(createdMonitor)
	if err != nil {
		log.Err(err).Msg("Error marshalling the monitor")
		http.Error(w, "Marshelling Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(createdMonitorJson)
}

// func (app *Application) getMonitorHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "Get Monitor")
// }

// func (app *Application) updateMonitorHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "Update Monitor")
// }

// func (app *Application) deleteMonitorHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "Delete Monitor")
// }
