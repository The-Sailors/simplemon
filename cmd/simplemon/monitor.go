package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/The-Sailors/simplemon/internal/data"
	"github.com/go-chi/httplog"
	"github.com/julienschmidt/httprouter"
)

func (app *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info().Msg("Starting Healthcheck Handler")
	fmt.Fprintln(w, "Jojo: is awesome!")
	fmt.Fprintln(w, "environment:", app.config.env)
}

func (app *Application) getAllMonitorsHandler(w http.ResponseWriter, r *http.Request) {
	log := httplog.LogEntry(r.Context())
	log.Info().Msg("Starting Get All Handler")
	//Get all the monitors from the database
	monitors, err := app.models.GetAll(r.Context(), log)
	if err != nil {
		log.Err(err).Msg("Error getting all the monitors")
		http.Error(w, "Error getting all the monitors", http.StatusInternalServerError)
		return
	}
	//Write the response
	monitorsJson, err := json.Marshal(monitors)
	if err != nil {
		log.Err(err).Msg("Error marshalling the monitor")
		http.Error(w, "Marshelling Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(monitorsJson)
}

func (app *Application) deleteMonitorHandler(w http.ResponseWriter, r *http.Request) {
	log := httplog.LogEntry(r.Context())
	log.Info().Msg("Starting Delete Handler")
	//get the id from the url
	monitorID := httprouter.ParamsFromContext(r.Context()).ByName("id")
	if monitorID == "" {
		log.Err(nil).Msg("Monitor id is required")
		http.Error(w, "Monitor id is required", http.StatusBadRequest)
		return
	}
	//convert the id to int
	monitorIDInt, err := strconv.Atoi(monitorID)
	if err != nil {
		log.Err(err).Msg("Error converting id to int")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//Verify if the monitor exists
	_, err = app.models.GetById(r.Context(), int64(monitorIDInt), log)
	if err != nil {
		if err.Error() == data.ErrMonitorNotFound.Error() {
			log.Warn().Msg("Monitor not found")
			http.Error(w, "Was not possible to delete the monitor because it not exists", http.StatusNotFound)
			return
		} else {
			log.Err(err).Msg("Error getting the monitor")
			http.Error(w, "Error getting the monitor", http.StatusInternalServerError)
			return
		}
	}
	//Delete the monitor
	err = app.models.Delete(r.Context(), int64(monitorIDInt), log)
	if err != nil {
		log.Err(err).Msg("Error deleting the monitor")
		http.Error(w, "Error deleting the monitor", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
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

		//verify if the error is a unique constraint violation
		if err.Error() == data.ErrUniqueConstraintViolation.Error() {
			log.Warn().Msg("Monitor already exists")
			http.Error(w, "Monitor already exists", http.StatusConflict)
			return
		} else {
			log.Err(err).Msg("Error creating the monitor")
			http.Error(w, "Error creating the monitor", http.StatusInternalServerError)
		}
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

func (app *Application) getMonitorHandler(w http.ResponseWriter, r *http.Request) {
	log := httplog.LogEntry(r.Context())
	log.Info().Msg("Starting Get Monitor Handler")
	params := httprouter.ParamsFromContext(r.Context())

	monitorID := params.ByName("id")

	if monitorID == "" {
		log.Err(nil).Msg("Monitor id is required")
		http.Error(w, "Monitor id is required", http.StatusBadRequest)
		return
	}
	//convert string to int64
	monitorIDInt, err := strconv.Atoi(monitorID)
	if err != nil {
		log.Err(err).Msgf("Error converting the monitor id: %s to int", monitorID)

		http.Error(w, "Invalid integer parameters", http.StatusBadRequest)
		return
	}

	//Get the monitor from the database
	monitor, err := app.models.GetById(r.Context(), int64(monitorIDInt), log)
	if err != nil {
		if err.Error() == data.ErrMonitorNotFound.Error() {
			log.Warn().Msg("Monitor not found")
			http.Error(w, "Monitor not found", http.StatusNotFound)
			return
		} else {

			log.Err(err).Msg("Error getting the monitor")
			http.Error(w, "Error getting the monitor", http.StatusInternalServerError)
			return
		}
	}
	monitorJson, err := json.Marshal(monitor)
	if err != nil {
		log.Err(err).Msg("Error marshalling the monitor")
		http.Error(w, "Marshelling Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(monitorJson)
}

// func (app *Application) updateMonitorHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "Update Monitor")
// }

// func (app *Application) deleteMonitorHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "Delete Monitor")
// }
