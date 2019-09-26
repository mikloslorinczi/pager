package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/mikloslorinczi/pager/client"
	"github.com/mikloslorinczi/pager/model"
	"github.com/mikloslorinczi/pager/rules"

	"github.com/spf13/cobra"
)

// oncallHoursCmd represents the user list command
var oncallHoursCmd = &cobra.Command{
	Use:     "hours",
	Aliases: []string{"h", "hour"},
	Short:   "Get infromation about on-call hours",
	Example: `Get Sandor's hours in thes month
pager -t PD_TOKEN oncall hours --user Sandor

Get Miklós's hours from last month
pager -t PD_TOKEN oncall hours --user Miklós --last

Get Viktor's hours since 2019.07.01 until 2019.09.30
pager -t PD_TOKEN oncall hours --user Viktor --since 2019.07.01 --until 2019.09.30

Get the Infrastructure team's hours from last month
pager -t PD_TOKEN oncall hours --team Infrastructure --last
`,
	Run: func(cmd *cobra.Command, args []string) {
		listOncallHours()
	},
}

func init() {
	oncallHoursCmd.Flags().String("user", "", "List only users with this name or ID")
	oncallHoursCmd.Flags().String("team", "", "List only users with this team-name or team-ID")
	oncallHoursCmd.Flags().Bool("last", false, "List on-calls from the last month")
	oncallHoursCmd.Flags().String("since", "", "List on-calls since this date e.g.: 2019.08.01")
	oncallHoursCmd.Flags().String("until", "", "List on-calls until this date e.g.: 2019.09.01")
	oncallCmd.AddCommand(oncallHoursCmd)
}

func listOncallHours() {
	req := client.GET()
	req.Path("/oncalls")
	if viper.GetBool("last") {
		viper.Set("since", rules.FormatDate(
			rules.BeginningOfMonth(time.Now().AddDate(0, -1, 0))))
		viper.Set("until", rules.FormatDate(
			rules.EndOfMonth(time.Now().AddDate(0, -1, 0))))
	}
	if since := viper.GetString("since"); since != "" {
		req.AddQuery("since", since)
	}
	if until := viper.GetString("until"); until != "" {
		req.AddQuery("until", until)
	}
	if user := viper.GetString("user"); user != "" {
		if userID := getUserID(user); userID != "" {
			req.AddQuery("user_ids[]", userID)
		}
	}
	res, err := req.Do()
	if err != nil {
		log.WithError(err).Fatal("HTTP Client failed to GET Response from Pager Futy API")
	}
	if !res.Ok {
		log.Fatalf("HTTP Error. Code: %d Body: %s", res.StatusCode, res.String())
	}
	var oncallsResp model.OncallsResponse
	if err := res.JSON(&oncallsResp); err != nil {
		log.WithError(err).Fatal("Failed to JSON Parse Response")
	}
	sumOnWork, sumOffWork := 0, 0
	fmt.Println()
	fmt.Println("---------------------------------------------------------")
	fmt.Println(" On-Call Hours")
	fmt.Println("---------------------------------------------------------")
	for _, oncall := range oncallsResp.Oncalls {
		fmt.Printf("User: %s\n", oncall.User.Summary)
		fmt.Printf("Escalation Policy: %s\n", oncall.EscalationPolicy.Summary)
		fmt.Printf("Escalation Level: %d\n", oncall.EscalationLevel)
		if schedule := oncall.Schedule.Summary; len(schedule) > 0 {
			fmt.Printf("Schedule: %s\n", oncall.Schedule.Summary)
			fmt.Printf("Start: %s - End: %s\n", oncall.Start, oncall.End)
			onWorkHours, offWorkHours := calculateWorkHours(oncall.Start, oncall.End, time.UTC)
			fmt.Printf("This covered %d OnWorkHours and %d OffWorkHours\n", onWorkHours, offWorkHours)
			sumOnWork += onWorkHours
			sumOffWork += offWorkHours
			fmt.Println()
		}
	}
	fmt.Println("---------------------------------------------------------")
	fmt.Printf("Sum of OnWorkHours: %d   Sum of OffWorkHours %d\n", sumOnWork, sumOffWork)
	fmt.Println("---------------------------------------------------------")
}

func getUserID(query string) string {
	req := client.GET()
	req.Path("/users")
	res, err := req.Do()
	if err != nil {
		log.WithError(err).Fatal("HTTP Client failed to GET Response from PagerDuty API")
	}
	var usersResp model.UsersResponse
	if err := res.JSON(&usersResp); err != nil {
		log.WithError(err).Fatal("Failed to JSON Parse Response")
	}
	for _, user := range usersResp.Users {
		if query == user.ID {
			return query
		}
		if strings.Contains(user.Name, query) {
			return user.ID
		}
	}
	return ""
}

func calculateWorkHours(start, end time.Time, location *time.Location) (int, int) {
	onWork, offWork := 0, 0
	elapsed := end.Sub(start)
	for offSet := time.Duration(0); offSet <= elapsed; offSet += time.Hour {
		// Check every hour (as a timestamp) from start to end
		sample := start.Add(offSet)
		if rules.IsWorkHour(sample, location) {
			onWork++
		} else {
			offWork++
		}
	}
	return onWork, offWork
}
