package models

// ChargeSettings represents charging configuration
type ChargeSettings struct {
	ID                     string  `json:"id"`
	CalculatedDeadline     string  `json:"calculatedDeadline"`
	Capacity               float64 `json:"capacity"`
	Deadline               string  `json:"deadline"`
	IsSmartChargingEnabled bool    `json:"isSmartChargingEnabled"`
	IsSolarChargingEnabled bool    `json:"isSolarChargingEnabled"`
	MaxChargeLimit         float64 `json:"maxChargeLimit"`
	MinChargeLimit         float64 `json:"minChargeLimit"`
	InitialCharge          float64 `json:"initialCharge"`
}

// ChargeState represents current charging state
type ChargeState struct {
	BatteryCapacity     float64 `json:"batteryCapacity"`
	BatteryLevel        float64 `json:"batteryLevel"`
	ChargeLimit         float64 `json:"chargeLimit"`
	ChargeRate          float64 `json:"chargeRate"`
	ChargeTimeRemaining float64 `json:"chargeTimeRemaining"`
	IsCharging          bool    `json:"isCharging"`
	IsFullyCharged      bool    `json:"isFullyCharged"`
	IsPluggedIn         bool    `json:"isPluggedIn"`
	LastUpdated         string  `json:"lastUpdated"`
	PowerDeliveryState  string  `json:"powerDeliveryState"`
	Range               float64 `json:"range"`
}

// DeviceInfo represents device information
type DeviceInfo struct {
	Brand string `json:"brand"`
	Model string `json:"model"`
	Year  int    `json:"year"`
	VIN   string `json:"vin,omitempty"`
}

// Intervention represents a device intervention/alert
type Intervention struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// EnodeCharger represents a smart charger
type EnodeCharger struct {
	ID             string          `json:"id"`
	CanSmartCharge bool            `json:"canSmartCharge"`
	ChargeSettings *ChargeSettings `json:"chargeSettings"`
	ChargeState    *ChargeState    `json:"chargeState"`
	Information    *DeviceInfo     `json:"information"`
	Interventions  []Intervention  `json:"interventions"`
	IsReachable    bool            `json:"isReachable"`
	LastSeen       string          `json:"lastSeen"`
}

// EnodeChargersResponse represents the API response
type EnodeChargersResponse struct {
	EnodeChargers []EnodeCharger `json:"enodeChargers"`
}

// EnodeVehicle represents a smart vehicle
type EnodeVehicle struct {
	ID             string          `json:"id"`
	CanSmartCharge bool            `json:"canSmartCharge"`
	ChargeSettings *ChargeSettings `json:"chargeSettings"`
	ChargeState    *ChargeState    `json:"chargeState"`
	Information    *DeviceInfo     `json:"information"`
	Interventions  []Intervention  `json:"interventions"`
	IsReachable    bool            `json:"isReachable"`
	LastSeen       string          `json:"lastSeen"`
}

// EnodeVehiclesResponse represents the API response
type EnodeVehiclesResponse struct {
	EnodeVehicles []EnodeVehicle `json:"enodeVehicles"`
}

// SmartBattery represents a smart battery
type SmartBattery struct {
	ID                string  `json:"id"`
	Brand             string  `json:"brand"`
	Capacity          float64 `json:"capacity"`
	MaxChargePower    float64 `json:"maxChargePower"`
	MaxDischargePower float64 `json:"maxDischargePower"`
	Provider          string  `json:"provider"`
	ExternalReference string  `json:"externalReference"`
	CreatedAt         string  `json:"createdAt"`
	UpdatedAt         string  `json:"updatedAt"`
}

// SmartBatteriesResponse represents the API response
type SmartBatteriesResponse struct {
	SmartBatteries []SmartBattery `json:"smartBatteries"`
}

// BatterySettings represents battery settings
type BatterySettings struct {
	BatteryMode                  string `json:"batteryMode"`
	ImbalanceTradingStrategy     string `json:"imbalanceTradingStrategy"`
	SelfConsumptionTradingAllowed bool   `json:"selfConsumptionTradingAllowed"`
}

// SmartBatteryDetails represents detailed battery info
type SmartBatteryDetails struct {
	ID       string           `json:"id"`
	Brand    string           `json:"brand"`
	Capacity float64          `json:"capacity"`
	Settings *BatterySettings `json:"settings"`
}

// SmartBatterySummary represents battery summary
type SmartBatterySummary struct {
	LastKnownStateOfCharge float64 `json:"lastKnownStateOfCharge"`
	LastKnownStatus        string  `json:"lastKnownStatus"`
	LastUpdate             string  `json:"lastUpdate"`
	TotalResult            float64 `json:"totalResult"`
}

// SmartBatteryDetailsResponse represents the API response
type SmartBatteryDetailsResponse struct {
	SmartBattery        *SmartBatteryDetails `json:"smartBattery"`
	SmartBatterySummary *SmartBatterySummary `json:"smartBatterySummary"`
}

// BatterySession represents a single battery session
type BatterySession struct {
	Date             string  `json:"date"`
	Result           float64 `json:"result"`
	CumulativeResult float64 `json:"cumulativeResult"`
	Status           string  `json:"status"`
	TradeIndex       float64 `json:"tradeIndex"`
}

// SmartBatterySessions represents battery sessions data
type SmartBatterySessions struct {
	DeviceID              string           `json:"deviceId"`
	FairUsePolicyVerified bool             `json:"fairUsePolicyVerified"`
	PeriodStartDate       string           `json:"periodStartDate"`
	PeriodEndDate         string           `json:"periodEndDate"`
	PeriodTotalResult     float64          `json:"periodTotalResult"`
	PeriodEpexResult      float64          `json:"periodEpexResult"`
	PeriodImbalanceResult float64          `json:"periodImbalanceResult"`
	PeriodTradingResult   float64          `json:"periodTradingResult"`
	PeriodFrankSlim       float64          `json:"periodFrankSlim"`
	PeriodTradeIndex      float64          `json:"periodTradeIndex"`
	Sessions              []BatterySession `json:"sessions"`
}

// SmartBatterySessionsResponse represents the API response
type SmartBatterySessionsResponse struct {
	SmartBatterySessions *SmartBatterySessions `json:"smartBatterySessions"`
}
