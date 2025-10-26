package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/fogleman/gg"
	"github.com/mightyzeus/stage-two/internal/models"
)

func fetchJSON(url string, target interface{}) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return err
	}
	return json.NewDecoder(resp.Body).Decode(&target)

}

func (app *application) GetAllCountries() ([]models.CountryApiResponse, error) {

	var countries []models.CountryApiResponse
	if err := fetchJSON("https://restcountries.com/v2/all?fields=name,capital,region,population,flag,currencies", &countries); err != nil {
		return nil, err
	}

	return countries, nil
}

func (app *application) GetExchangeRates() (*models.ExchangeRateResponse, error) {

	var rates models.ExchangeRateResponse
	if err := fetchJSON("https://open.er-api.com/v6/latest/USD", &rates); err != nil {
		return nil, err
	}

	return &rates, nil
}

func CalculateEstimatedGDP(population int64, exchangeRate float64) float64 {
	randomFactor := float64(rand.Intn(1001) + 1000)

	return (float64(population) * randomFactor) / exchangeRate

}

func GenerateSummaryImage(countries []models.Country, outputPath string) error {
	const W = 600
	const H = 400

	dc := gg.NewContext(W, H)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Text settings
	dc.SetRGB(0, 0, 0) // black color
	if err := dc.LoadFontFace("/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf", 20); err != nil {
		log.Println("Failed to load font:", err)
	}

	y := 40

	// Total number of countries
	total := len(countries)
	dc.DrawStringAnchored(fmt.Sprintf("Total countries: %d", total), W/2, float64(y), 0.5, 0.5)
	y += 40

	// Sort countries by EstimatedGdp descending
	sort.Slice(countries, func(i, j int) bool {
		if countries[i].EstimatedGdp == nil {
			return false
		}
		if countries[j].EstimatedGdp == nil {
			return true
		}
		return *countries[i].EstimatedGdp > *countries[j].EstimatedGdp
	})

	// Top 5 countries by GDP
	dc.DrawStringAnchored("Top 5 Countries by Estimated GDP:", W/2, float64(y), 0.5, 0.5)
	y += 30
	for i := 0; i < 5 && i < len(countries); i++ {
		c := countries[i]
		gdp := 0.0
		if c.EstimatedGdp != nil {
			gdp = *c.EstimatedGdp
		}
		dc.DrawStringAnchored(fmt.Sprintf("%d. %s - %.2f", i+1, c.Name, gdp), W/2, float64(y), 0.5, 0.5)
		y += 25
	}

	// Last refreshed timestamp (take the max LastRefreshedAt)
	var lastRefreshed time.Time
	for _, c := range countries {
		if c.LastRefreshedAt.After(lastRefreshed) {
			lastRefreshed = c.LastRefreshedAt
		}
	}
	y += 20
	dc.DrawStringAnchored(fmt.Sprintf("Last refreshed: %s", lastRefreshed.Format(time.RFC1123)), W/2, float64(y), 0.5, 0.5)

	// Create cache directory if not exists
	if err := os.MkdirAll("cache", os.ModePerm); err != nil {
		return err
	}

	return dc.SavePNG(outputPath)
}
