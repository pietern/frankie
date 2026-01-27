package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// DefaultURL is the Frank Energie GraphQL API endpoint
	DefaultURL = "https://graphql.frankenergie.nl/"

	// DefaultTimeout is the default HTTP timeout
	DefaultTimeout = 30 * time.Second

	// ClientVersion is the version string for API headers
	ClientVersion = "4.13.3"
	ClientName    = "frank-app"
	ClientOS      = "ios/26.0.1"
)

// Client is the GraphQL client for Frank Energie API
type Client struct {
	httpClient *http.Client
	baseURL    string
	authToken  string
	country    string
}

// NewClient creates a new API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		baseURL: DefaultURL,
		country: "NL",
	}
}

// SetAuthToken sets the authentication token
func (c *Client) SetAuthToken(token string) {
	c.authToken = token
}

// SetCountry sets the country for API requests (NL or BE)
func (c *Client) SetCountry(country string) {
	c.country = country
}

// GraphQLRequest represents a GraphQL request
type GraphQLRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse represents a GraphQL response
type GraphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []GraphQLError  `json:"errors,omitempty"`
}

// Execute sends a GraphQL request and returns the raw response
func (c *Client) Execute(query, operationName string, variables map[string]interface{}) (*GraphQLResponse, error) {
	return c.ExecuteWithHeaders(query, operationName, variables, nil)
}

// ExecuteWithHeaders sends a GraphQL request with custom headers
func (c *Client) ExecuteWithHeaders(query, operationName string, variables map[string]interface{}, extraHeaders map[string]string) (*GraphQLResponse, error) {
	reqBody := GraphQLRequest{
		Query:         query,
		OperationName: operationName,
		Variables:     variables,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-graphql-client-version", ClientVersion)
	req.Header.Set("x-graphql-client-name", ClientName)
	req.Header.Set("x-graphql-client-os", ClientOS)
	req.Header.Set("skip-graphcdn", "1")

	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	}

	if c.country != "" && c.country != "NL" {
		req.Header.Set("x-country", c.country)
	}

	for k, v := range extraHeaders {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNetwork, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrAuthRequired
	}
	if resp.StatusCode == http.StatusForbidden {
		return nil, ErrForbidden
	}
	if resp.StatusCode == http.StatusBadRequest {
		return nil, ErrBadRequest
	}
	if resp.StatusCode >= 500 {
		return nil, ErrServerError
	}

	var gqlResp GraphQLResponse
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for GraphQL errors
	if len(gqlResp.Errors) > 0 {
		return &gqlResp, c.handleGraphQLErrors(gqlResp.Errors)
	}

	return &gqlResp, nil
}

// handleGraphQLErrors processes GraphQL errors and returns appropriate Go errors
func (c *Client) handleGraphQLErrors(errors []GraphQLError) error {
	if len(errors) == 0 {
		return nil
	}

	msg := errors[0].Message

	switch msg {
	case "user-error:password-invalid":
		return ErrInvalidCredentials
	case "user-error:auth-not-authorised":
		return ErrForbidden
	case "user-error:auth-required":
		return ErrAuthRequired
	case "user-error:smart-trading-not-enabled":
		return ErrSmartTradingNotEnabled
	case "user-error:smart-charging-not-enabled":
		return ErrSmartChargingNotEnabled
	case "request-error:request-not-supported-in-country":
		return ErrNotSupportedInCountry
	}

	return &APIError{Message: msg}
}

// Login authenticates with email and password
func (c *Client) Login(email, password string) (authToken, refreshToken string, err error) {
	variables := map[string]interface{}{
		"email":    email,
		"password": password,
	}

	resp, err := c.Execute(LoginMutation, "Login", variables)
	if err != nil {
		return "", "", err
	}

	var result struct {
		Login struct {
			AuthToken    string `json:"authToken"`
			RefreshToken string `json:"refreshToken"`
		} `json:"login"`
	}

	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return "", "", fmt.Errorf("failed to parse login response: %w", err)
	}

	return result.Login.AuthToken, result.Login.RefreshToken, nil
}

// RenewToken refreshes the authentication token
func (c *Client) RenewToken(authToken, refreshToken string) (newAuthToken, newRefreshToken string, err error) {
	variables := map[string]interface{}{
		"authToken":    authToken,
		"refreshToken": refreshToken,
	}

	resp, err := c.Execute(RenewTokenMutation, "RenewToken", variables)
	if err != nil {
		return "", "", err
	}

	var result struct {
		RenewToken struct {
			AuthToken    string `json:"authToken"`
			RefreshToken string `json:"refreshToken"`
		} `json:"renewToken"`
	}

	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return "", "", fmt.Errorf("failed to parse renew token response: %w", err)
	}

	return result.RenewToken.AuthToken, result.RenewToken.RefreshToken, nil
}
