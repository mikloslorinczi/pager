package cmd

import (
	"github.com/spf13/cobra"
)

// teamCmd represents the run command
var teamCmd = &cobra.Command{
	Use:     "team",
	Aliases: []string{"t", "teams", "tm", "te"},
	Short:   "Team Information",
	Long: `
Get Team information
	`,
	Example: "pager -t PD_TOKEN team list",
}

func init() {
	rootCmd.AddCommand(teamCmd)
}
