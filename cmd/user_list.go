package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/mikloslorinczi/pager/client"
	"github.com/mikloslorinczi/pager/model"
	"github.com/spf13/cobra"
)

// userListCmd represents the user list command
var userListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls", "ps"},
	Short:   "List users",
	Example: "pager -t PD_TOKEN user list --team Infra",
	Run: func(cmd *cobra.Command, args []string) {
		listUsers()
	},
}

func init() {
	userListCmd.Flags().String("team", "", "Only list users belonging to this team (team name, or team id)")
	userCmd.AddCommand(userListCmd)
}

func listUsers() {
	req := client.GET()
	req.Path("/users")
	res, err := req.Do()
	if err != nil {
		log.WithError(err).Fatal("HTTP Client failed to GET Response from PagerDuty API")
	}
	if !res.Ok {
		log.Fatalf("HTTP Error. Code: %d Body: %s", res.StatusCode, res.String())
	}
	var usersResp model.UsersResponse
	if err := res.JSON(&usersResp); err != nil {
		log.WithError(err).Fatal("Failed to JSON Parse Response")
	}
	for _, user := range usersResp.Users {
		if team := viper.GetString("team"); team != "" {
			if !userInTeam(user, viper.GetString("team")) {
				continue
			}
		}
		fmt.Printf("ID: %s Name: %s Teams: ", user.ID, user.Name)
		for _, team := range user.Teams {
			fmt.Printf("%s, ", team.Summary)
		}
		fmt.Println()
	}
}

func userInTeam(user model.User, team string) bool {
	if len(team) < 1 {
		return false
	}
	for _, t := range user.Teams {
		if strings.Contains(strings.ToLower(t.Summary), strings.ToLower(team)) || t.ID == team {
			return true
		}
	}
	return false
}
