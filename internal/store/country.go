package store

import (
	"context"
	"strings"

	"github.com/mightyzeus/stage-two/internal/models"
	"gorm.io/gorm"
)

type CountryStore struct {
	db *gorm.DB
}

func (c *CountryStore) CountryRefresh(ctx context.Context, countries []models.Country) error {
	if len(countries) == 0 {
		return nil
	}

	// lowercase the country name
	names := make([]string, len(countries))
	for i, country := range countries {
		names[i] = strings.ToLower(country.Name)
	}

	var existingCountries []models.Country
	if err := c.db.WithContext(ctx).
		Where("LOWER(name) IN ?", names).
		Find(&existingCountries).Error; err != nil {
		return err
	}

	// use map for the existing countries so that the query is faster...make sense that way instead of mapping or ranging
	existingMap := make(map[string]*models.Country)
	for i := range existingCountries {
		existingMap[strings.ToLower(existingCountries[i].Name)] = &existingCountries[i]
	}

	var toInsert []models.Country

	for _, country := range countries {
		lcName := strings.ToLower(country.Name)
		if existing, ok := existingMap[lcName]; ok {
			// Update existing records
			existing.Region = country.Region
			existing.Population = country.Population
			existing.CurrencyCode = country.CurrencyCode
			existing.ExchangeRate = country.ExchangeRate
			existing.EstimatedGdp = country.EstimatedGdp
			existing.FlagUrl = country.FlagUrl
			existing.LastRefreshedAt = country.LastRefreshedAt
		} else {
			// add new records if e no exits
			toInsert = append(toInsert, country)
		}
	}

	// Bulk add new countries
	if len(toInsert) > 0 {
		if err := c.db.WithContext(ctx).Create(&toInsert).Error; err != nil {
			return err
		}
	}

	// Bulk update existing countries
	for _, existing := range existingMap {
		if err := c.db.WithContext(ctx).Save(existing).Error; err != nil {
			return err
		}
	}

	return nil
}

func (c *CountryStore) GetCountryByName(ctx context.Context, name string) (*models.Country, error) {
	var country models.Country
	err := c.db.WithContext(ctx).
		Where("LOWER(name) = ?", strings.ToLower(name)).
		First(&country).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &country, nil
}

func (c *CountryStore) DeleteCountryByName(ctx context.Context, name string) (bool, error) {
	res := c.db.WithContext(ctx).
		Where("LOWER(name) = ?", strings.ToLower(name)).
		Delete(&models.Country{})

	if res.Error != nil {
		return false, res.Error
	}

	if res.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (c *CountryStore) GetAllCountries(ctx context.Context, region, currency, sort string) ([]*models.Country, int64, error) {
	var countries []*models.Country
	var total int64

	db := c.db.WithContext(ctx).Model(&models.Country{})

	if region != "" {
		db = db.Where("LOWER(region) = ?", strings.ToLower(region))
	}
	if currency != "" {
		db = db.Where("currency_code = ?", strings.ToUpper(currency))
	}

	switch sort {
	case "gdp_desc":
		db = db.Order("estimated_gdp DESC")
	case "gdp_asc":
		db = db.Order("estimated_gdp ASC")
	case "population_desc":
		db = db.Order("population DESC")
	case "population_asc":
		db = db.Order("population ASC")
	default:
		db = db.Order("name ASC")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Find(&countries).Error; err != nil {
		return nil, 0, err
	}

	return countries, total, nil
}
