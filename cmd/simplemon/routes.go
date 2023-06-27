package main

import (
	"net/http"

	"github.com/go-chi/httplog"
	"github.com/go-openapi/runtime/middleware"
	"github.com/julienschmidt/httprouter"
)

func (app *Application) routes() http.Handler {
	httpLogMiddleware := httplog.RequestLogger(app.logger)

	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/v1/monitors", addMiddleware(app.createMonitorHandler, httpLogMiddleware))
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", addMiddleware(app.healthcheckHandler, httpLogMiddleware))
	router.HandlerFunc(http.MethodGet, "/v1/monitors/:id", addMiddleware(app.getMonitorHandler, httpLogMiddleware))
	router.HandlerFunc(http.MethodDelete, "/v1/monitors/:id", addMiddleware(app.deleteMonitorHandler, httpLogMiddleware))
	opts := middleware.SwaggerUIOpts{SpecURL: "openapi.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	router.Handler(http.MethodGet, "/docs", sh)
	router.Handler(http.MethodGet, "/openapi.yaml", http.FileServer(http.Dir("./")))
	return router
}

func addMiddleware(handler http.HandlerFunc, middleware func(next http.Handler) http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Apply the middleware to the handler
		middleware(handler).ServeHTTP(w, r)
	}
}
