package metrics

import (
	"fmt"
	"net/http"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/rs/zerolog"

	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/loophole/cli/internal/app/metrics/handlers"
	"github.com/rs/zerolog/log"
)

func Serve() {
	metricsLogger := log.
		With().
		Timestamp().
		Str("component", "metrics").
		Logger().
		Level(zerolog.InfoLevel)
	appBox, err := rice.FindBox("static/dashboard/build")
	if err != nil {
		metricsLogger.Error().Err(err).Msg("Failed to find frontend files")
	}
	r := mux.NewRouter()
	r.HandleFunc("/api", handlers.EventsHandler)
	r.HandleFunc("/api/current", handlers.CurrentSiteHandler).
		Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/events", handlers.EventsHandler).
		Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/events/{siteId}", handlers.SiteEventsHandler).
		Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/metrics", handlers.MetricsHandler).
		Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/metrics/{siteId}", handlers.SiteMetricsHandler).
		Methods(http.MethodGet, http.MethodOptions)

	r.PathPrefix("/").Handler(http.FileServer(appBox.HTTPBox()))
	r.Use(mux.CORSMethodMiddleware(r))

	srv := &http.Server{
		Handler:      ghandlers.LoggingHandler(metricsLogger, r),
		Addr:         ":12345",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	metricsLogger.Info().Msg(fmt.Sprintf("Metrics server listening on http://localhost%s", srv.Addr))
	if err := srv.ListenAndServe(); err != nil {
		metricsLogger.Error().Err(err).Msg("Failed to run metrics server")
	}
}

func serveAppHandler(appBox *rice.Box) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexFile, err := appBox.Open("index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.ServeContent(w, r, "index.html", time.Time{}, indexFile)
	}
}
