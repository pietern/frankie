package models

import "time"

// Price represents a single price entry
type Price struct {
	From               time.Time `json:"from"`
	Till               time.Time `json:"till"`
	Resolution         string    `json:"resolution"`
	MarketPrice        float64   `json:"marketPrice"`
	MarketPriceTax     float64   `json:"marketPriceTax"`
	SourcingMarkupPrice float64  `json:"sourcingMarkupPrice"`
	EnergyTaxPrice     float64   `json:"energyTaxPrice"`
	MarketPricePlus    float64   `json:"marketPricePlus"`
	AllInPrice         float64   `json:"allInPrice"`
	PerUnit            string    `json:"perUnit"`
}

// TotalPrice returns the total price (market + tax + markup + energy tax)
func (p *Price) TotalPrice() float64 {
	return p.MarketPrice + p.MarketPriceTax + p.SourcingMarkupPrice + p.EnergyTaxPrice
}

// AveragePrice represents average price information
type AveragePrice struct {
	AverageMarketPrice     float64 `json:"averageMarketPrice"`
	AverageMarketPricePlus float64 `json:"averageMarketPricePlus"`
	AverageAllInPrice      float64 `json:"averageAllInPrice"`
	PerUnit                string  `json:"perUnit"`
	IsWeighted             bool    `json:"isWeighted"`
}

// MarketPrices represents market prices response
type MarketPrices struct {
	AverageElectricityPrices *AveragePrice `json:"averageElectricityPrices"`
	ElectricityPrices        []Price       `json:"electricityPrices"`
	GasPrices                []Price       `json:"gasPrices"`
}

// MarketPricesResponse represents the API response for market prices
type MarketPricesResponse struct {
	MarketPrices *MarketPrices `json:"marketPrices"`
}

// CustomerMarketPricesResponse represents the API response for customer prices
type CustomerMarketPricesResponse struct {
	CustomerMarketPrices *MarketPrices `json:"customerMarketPrices"`
}
