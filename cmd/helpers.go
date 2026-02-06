package cmd

import (
	"fmt"

	"github.com/pietern/frankie/internal/api"
	"github.com/pietern/frankie/internal/auth"
)

// newAuthenticatedClient creates an API client and ensures the user is logged in.
func newAuthenticatedClient() (*api.Client, error) {
	client := api.NewClient()
	manager := auth.NewManager(client)
	if err := manager.EnsureAuthenticated(); err != nil {
		return nil, fmt.Errorf("not logged in: %w", err)
	}
	return client, nil
}
