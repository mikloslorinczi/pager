package cmd

import (
	"fmt"

	"github.com/mikloslorinczi/pager/client"

	"github.com/spf13/cobra"
)

// teamListCmd represents the user list command
var teamListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls", "ps"},
	Short:   "List teams",
	Example: "pager -t PD_TOKEN team list",
	Run: func(cmd *cobra.Command, args []string) {
		listTeams()
	},
}

func init() {
	teamCmd.AddCommand(teamListCmd)
}

func listTeams() {
	teamsResp, err := client.GetTeamsList()
	if err != nil {
		log.WithError(err).Fatal("Failed to get the list of teams")
	}
	for _, team := range teamsResp.Teams {
		fmt.Printf("ID: %s Name: %s\n", team.ID, team.Name)
	}
}
