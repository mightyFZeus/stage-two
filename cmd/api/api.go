package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mightyzeus/stage-two/internal/store"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr   string
	apiUrl string
	db     dbConfig
}

type dbConfig struct {
	dbAddr       string
	maxOpenConns int
	maxIdleTime  string
	maxIdleConns int
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Post("/countries/refresh", app.RefreshCountriesHandler)
	r.Get("/countries/image", app.ServeSummaryImage)
	r.Route("/countries", func(r chi.Router) {
		r.Get("/", app.GetCountriesHandler)
		r.Get("/{name}", app.GetCountryByName)
	})
	r.Delete("/countries/{name}", app.DeleteCountryByName)
	r.Get("/status", app.GetStatus)

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	return srv.ListenAndServe()
}
