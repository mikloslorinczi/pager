package cmd

import (
	"github.com/spf13/cobra"
)

// userCmd represents the run command
var userCmd = &cobra.Command{
	Use:     "user",
	Aliases: []string{"u", "us", "users"},
	Short:   "User Infromation",
	Long: `
Get User information
	`,
	Example: "pager -t PD_TOKEN user list",
}

func init() {
	rootCmd.AddCommand(userCmd)
}
