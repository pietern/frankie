package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pietern/frankie/internal/auth"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear stored credentials",
	Long:  `Remove stored authentication credentials from the local machine.`,
	RunE:  runLogout,
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}

func runLogout(cmd *cobra.Command, args []string) error {
	if !auth.CredentialsExist() {
		fmt.Println("Not logged in")
		return nil
	}

	if err := auth.DeleteCredentials(); err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	fmt.Println("Logged out successfully")
	return nil
}
