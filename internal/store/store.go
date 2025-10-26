package store

import (
	"context"

	"github.com/mightyzeus/stage-two/internal/models"
	"gorm.io/gorm"
)

type Storage struct {
	Country interface {
		CountryRefresh(ctx context.Context, countries []models.Country) error
		GetCountryByName(ctx context.Context, value string) (*models.Country, error)
		DeleteCountryByName(ctx context.Context, name string) (bool, error)
		GetAllCountries(ctx context.Context, region, currency, sort string) ([]*models.Country, int64, error)
	}
}

func NewStorage(db *gorm.DB) Storage {
	return Storage{
		Country: &CountryStore{db},
	}
}
