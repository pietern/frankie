package models

import "time"

// User represents user information from the API
type User struct {
	ID                     string           `json:"id"`
	Email                  string           `json:"email"`
	CountryCode            string           `json:"countryCode"`
	AdvancedPaymentAmount  float64          `json:"advancedPaymentAmount"`
	TreesCount             int              `json:"treesCount"`
	HasInviteLink          bool             `json:"hasInviteLink"`
	HasCO2Compensation     bool             `json:"hasCO2Compensation"`
	CreatedAt              time.Time        `json:"createdAt"`
	UpdatedAt              time.Time        `json:"updatedAt"`
	ExternalDetails        *ExternalDetails `json:"externalDetails"`
	SmartCharging          *SmartCharging   `json:"smartCharging"`
	SmartTrading           *SmartTrading    `json:"smartTrading"`
	WebsiteURL             string           `json:"websiteUrl"`
	CustomerSupportEmail   string           `json:"customerSupportEmail"`
	Reference              string           `json:"reference"`
	Connections            []Connection     `json:"connections"`
}

// ExternalDetails contains external user details
type ExternalDetails struct {
	Reference string   `json:"reference"`
	Person    *Person  `json:"person"`
	Contact   *Contact `json:"contact"`
	Address   *Address `json:"address"`
}

// Person represents name information
type Person struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// Contact represents contact information
type Contact struct {
	EmailAddress string `json:"emailAddress"`
	PhoneNumber  string `json:"phoneNumber"`
	MobileNumber string `json:"mobileNumber"`
}

// Address represents address information
type Address struct {
	AddressFormatted    []string `json:"addressFormatted"`
	Street              string   `json:"street"`
	HouseNumber         string   `json:"houseNumber"`
	HouseNumberAddition string   `json:"houseNumberAddition"`
	ZipCode             string   `json:"zipCode"`
	City                string   `json:"city"`
}

// FormattedAddress returns the address as a single string
func (a *Address) FormattedAddress() string {
	if len(a.AddressFormatted) > 0 {
		return a.AddressFormatted[0]
	}
	return ""
}

// SmartCharging represents smart charging status
type SmartCharging struct {
	IsActivated          bool   `json:"isActivated"`
	Provider             string `json:"provider"`
	IsAvailableInCountry bool   `json:"isAvailableInCountry"`
}

// SmartTrading represents smart trading status
type SmartTrading struct {
	IsActivated          bool `json:"isActivated"`
	IsAvailableInCountry bool `json:"isAvailableInCountry"`
}

// Site represents a user site
type Site struct {
	Address               *SiteAddress `json:"address"`
	AddressHasMultipleSites bool       `json:"addressHasMultipleSites"`
	DeliveryEndDate       string       `json:"deliveryEndDate"`
	DeliveryStartDate     string       `json:"deliveryStartDate"`
	FirstMeterReadingDate string       `json:"firstMeterReadingDate"`
	LastMeterReadingDate  string       `json:"lastMeterReadingDate"`
	PropositionType       string       `json:"propositionType"`
	Reference             string       `json:"reference"`
	Segments              []string     `json:"segments"`
	Status                string       `json:"status"`
}

// SiteAddress represents a site's address
type SiteAddress struct {
	AddressFormatted []string `json:"addressFormatted"`
}

// FormattedAddress returns the address as a single string
func (a *SiteAddress) FormattedAddress() string {
	if len(a.AddressFormatted) > 0 {
		return a.AddressFormatted[0]
	}
	return ""
}

// Connection represents an energy connection (electricity or gas)
type Connection struct {
	ID                    string                    `json:"id"`
	ConnectionID          string                    `json:"connectionId"`
	EAN                   string                    `json:"EAN"`
	Segment               string                    `json:"segment"`
	Status                string                    `json:"status"`
	ContractStatus        string                    `json:"contractStatus"`
	EstimatedFeedIn       float64                   `json:"estimatedFeedIn"`
	FirstMeterReadingDate string                    `json:"firstMeterReadingDate"`
	LastMeterReadingDate  string                    `json:"lastMeterReadingDate"`
	MeterType             string                    `json:"meterType"`
	ExternalDetails       *ConnectionExternalDetails `json:"externalDetails"`
}

// ConnectionExternalDetails contains external connection details
type ConnectionExternalDetails struct {
	GridOperator string           `json:"gridOperator"`
	Address      *Address         `json:"address"`
	Contract     *ContractDetails `json:"contract"`
}

// ContractDetails contains contract information
type ContractDetails struct {
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	ContractType  string `json:"contractType"`
	ProductName   string `json:"productName"`
	TariffChartID string `json:"tariffChartId"`
}

// MeResponse represents the response from the Me query
type MeResponse struct {
	Me User `json:"me"`
}

// UserSitesResponse represents the response from the UserSites query
type UserSitesResponse struct {
	UserSites []Site `json:"userSites"`
}
