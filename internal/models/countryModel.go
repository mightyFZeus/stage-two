package models

import "time"

type Country struct {
	ID              string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name            string    `json:"name" validate:"required"`
	Region          string    `json:"region" validate:"required"`
	Population      int       `json:"population" validate:"required"`
	CurrencyCode    *string   `json:"currency_code" validate:"required"`
	ExchangeRate    *float64  `json:"exchange_rate" validate:"required"`
	EstimatedGdp    *float64  `json:"estimated_gdp" validate:"required"`
	FlagUrl         string    `json:"flag_url" validate:"required"`
	LastRefreshedAt time.Time `json:"last_refreshed_at" validate:"required"`
}
type Currency struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}
type CountryApiResponse struct {
	Name        string     `json:"name"`
	Capital     string     `json:"capital"`
	Region      string     `json:"region"`
	Population  int        `json:"population"`
	Flag        string     `json:"flag"`
	Independent bool       `json:"independent"`
	Currencies  []Currency `json:"currencies"`
}

type ExchangeRateResponse struct {
	Result            string             `json:"result"`
	Provider          string             `json:"provider"`
	Documentation     string             `json:"documentation"`
	TermsOfUse        string             `json:"terms_of_use"`
	TimeLastUpdate    int64              `json:"time_last_update_unix"`
	TimeLastUpdateUTC string             `json:"time_last_update_utc"`
	TimeNextUpdate    int64              `json:"time_next_update_unix"`
	TimeNextUpdateUTC string             `json:"time_next_update_utc"`
	TimeEOLUnix       int64              `json:"time_eol_unix"`
	BaseCode          string             `json:"base_code"`
	Rates             map[string]float64 `json:"rates"`
}
