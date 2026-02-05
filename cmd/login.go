package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/pietern/frankie/internal/api"
	"github.com/pietern/frankie/internal/auth"
)

var (
	loginEmail    string
	loginPassword string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Frank Energie",
	Long:  `Login to Frank Energie with your email and password.`,
	RunE:  runLogin,
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&loginEmail, "email", "e", "", "email address")
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "password")
}

func runLogin(cmd *cobra.Command, args []string) error {
	email := loginEmail
	password := loginPassword

	// If credentials not provided via flags, use interactive form
	if email == "" || password == "" {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Email").
					Value(&email).
					Placeholder("your@email.com"),
				huh.NewInput().
					Title("Password").
					Value(&password).
					EchoMode(huh.EchoModePassword),
			),
		)

		if err := form.Run(); err != nil {
			return fmt.Errorf("form cancelled: %w", err)
		}
	}

	if email == "" || password == "" {
		return fmt.Errorf("email and password are required")
	}

	client := api.NewClient()
	manager := auth.NewManager(client)

	err := manager.Login(email, password)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Login failed:", err)
		return err
	}

	fmt.Println("Login successful!")

	return nil
}
