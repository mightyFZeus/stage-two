package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mightyzeus/stage-two/internal/models"
)

func (app *application) RefreshCountriesHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	countries, err := app.GetAllCountries()
	if err != nil {
		log.Println(err)
		writeJSONError(w, http.StatusBadGateway, "External data source unavailable", "Could not fetch data from RestCountry Api")

		return
	}

	rates, err := app.GetExchangeRates()
	if err != nil {
		writeJSONError(w, http.StatusBadGateway, "External data source unavailable", "Could not fetch data from Exchange Rate Api")
		return
	}

	var matched []models.Country
	now := time.Now()

	for _, country := range countries {
		var code *string
		var rate *float64
		var gdp *float64
		zero := float64(0)

		if len(country.Currencies) > 0 {
			firstCurrency := country.Currencies[0].Code
			code = &firstCurrency

			if val, ok := rates.Rates[firstCurrency]; ok {
				rate = &val
				estGDP := CalculateEstimatedGDP(int64(country.Population), val)
				gdp = &estGDP
			} else {
				// Currency code not found
				gdp = &zero
			}
		} else {
			// No currencies array
			gdp = &zero
		}

		c := models.Country{
			ID:              uuid.New().String(),
			Name:            country.Name,
			Region:          country.Region,
			Population:      country.Population,
			CurrencyCode:    code,
			ExchangeRate:    rate,
			EstimatedGdp:    gdp,
			FlagUrl:         country.Flag,
			LastRefreshedAt: now,
		}

		matched = append(matched, c)
	}

	if err := app.store.Country.CountryRefresh(ctx, matched); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := GenerateSummaryImage(matched, "cache/summary.png"); err != nil {
		log.Println("Failed to generate summary image:", err)
	}

	app.jsonResponse(w, http.StatusOK, map[string]string{
		"message": "countries refreshed successfully",
	})
}

func (app *application) ServeSummaryImage(w http.ResponseWriter, r *http.Request) {
	const filePath = "cache/summary.png"

	// Check if file exists
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			app.notFoundResponse(w, r, errors.New("summary image not found"))
			return
		}
		http.Error(w, "Error accessing summary image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	http.ServeFile(w, r, filePath)
}

func (app *application) GetCountryByName(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	str := strings.TrimSpace(strings.ToLower(chi.URLParam(r, "name")))
	if str == "" {
		writeJSONError(w, http.StatusBadRequest, "Validation failed", "country_name is required")

		return
	}

	existingWord, _ := app.store.Country.GetCountryByName(ctx, str)
	if existingWord == nil {
		app.notFoundResponse(w, r, errors.New("country not found"))
		return
	}

	app.jsonResponse(w, http.StatusOK, existingWord)
}
func (app *application) DeleteCountryByName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	name := strings.TrimSpace(strings.ToLower(chi.URLParam(r, "name")))
	if name == "" {
		writeJSONError(w, http.StatusBadRequest, "Validation failed", "country_name is required")
		return
	}

	deleted, err := app.store.Country.DeleteCountryByName(ctx, name)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if !deleted {
		writeJSONError(w, http.StatusNotFound, "Validation failed", "country does not exist")
		return
	}

	app.jsonResponse(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("country '%s' deleted successfully", name),
	})
}

func (app *application) GetCountriesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	region := r.URL.Query().Get("region")
	currency := r.URL.Query().Get("currency")
	sort := r.URL.Query().Get("sort")

	countries, _, err := app.store.Country.GetAllCountries(ctx, region, currency, sort)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, map[string]interface{}{
		"data": countries,
	})
}

func (app *application) GetStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	countries, total, err := app.store.Country.GetAllCountries(ctx, "", "", "")
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Find the most recent LastRefreshedAt
	var lastRefreshed time.Time
	for _, c := range countries {
		if c.LastRefreshedAt.After(lastRefreshed) {
			lastRefreshed = c.LastRefreshedAt
		}
	}

	app.jsonResponse(w, http.StatusOK, map[string]interface{}{
		"total_countries": total,
		"last_refresh":    lastRefreshed.Format(time.RFC3339),
	})
}
