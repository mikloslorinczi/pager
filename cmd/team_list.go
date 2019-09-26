package cmd

import (
	"fmt"

	"github.com/mikloslorinczi/pager/client"
	"github.com/mikloslorinczi/pager/model"
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
	req := client.GET()
	req.Path("/teams")
	res, err := req.Do()
	if err != nil {
		log.WithError(err).Fatal("HTTP Client failed to GET Response from Pager Futy API")
	}
	if !res.Ok {
		log.Fatalf("HTTP Error. Code: %d Body: %s", res.StatusCode, res.String())
	}
	var teamsResp model.TeamsResponse
	if err := res.JSON(&teamsResp); err != nil {
		log.WithError(err).Fatal("Failed to JSON Parse Response")
	}
	for _, team := range teamsResp.Teams {
		fmt.Printf("ID: %s Name: %s\n", team.ID, team.Name)
	}
}
