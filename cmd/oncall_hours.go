package cmd

import (
	"fmt"
	"time"

	"github.com/mikloslorinczi/pager/config"
	"github.com/mikloslorinczi/pager/client"
	"github.com/mikloslorinczi/pager/model"
	"github.com/mikloslorinczi/pager/rules"

	"github.com/spf13/viper"
	"github.com/spf13/cobra"
)

// oncallHoursCmd represents the user list command
var oncallHoursCmd = &cobra.Command{
	Use:     "hours",
	Aliases: []string{"h", "hour"},
	Short:   "Get information about on-call hours",
	Long: `Get infromation about on-call hours
This command calculates the offWork and onWork on-call hours
of everyone, a team or a single user

If you set the --user flag with an user's ID or name (or part of the name) you will
get only that person's data.

If you set the --team flag which also accepts an ID or a name (or part of the name)
all members of the team will be added to the query

The default timeframe of the query is the begining of this month and the time of execution
You can set the timeframe with the --since and --until flags in the 2019-09-10 format
By setting the --last flag the timeframe will be the begining and the end of the last month
`,
	Example: `Get Sandor's hours in this month
pager -t PD_TOKEN oncall hours --user Sandor

Get Miklós's hours from last month
pager -t PD_TOKEN oncall hours --user Miklós --last

Get Viktor's hours since 2019.07.01 until 2019.09.30
pager -t PD_TOKEN oncall hours --user Viktor --since 2019-07-01 --until 2019-09-30

Get the Infrastructure team's hours from last month
pager -t PD_TOKEN oncall hours --team Infrastructure --last
`,
	Run: func(cmd *cobra.Command, args []string) {
		listOncallHours()
	},
}

var since, until time.Time

func init() {
	oncallHoursCmd.Flags().String("user", "", "Add this user's on-calls (by username or ID) to the list")
	oncallHoursCmd.Flags().String("team", "", "Add all member's on-calls from this team (by teamname or ID) to the list")
	oncallHoursCmd.Flags().Bool("last", false, "List on-calls from the last month")
	oncallHoursCmd.Flags().String("since", "", "List on-calls since this date e.g.: 2019-08-01")
	oncallHoursCmd.Flags().String("until", "", "List on-calls until this date e.g.: 2019-09-01")
	oncallCmd.AddCommand(oncallHoursCmd)
}

// setSinceUntil sets the since and until times in between we list the on-calls
// By default since is the beginning of this month and until is now
// If --last is set since will be the beginning of the last month and until will be its end
func setSinceUntil() {
	if viper.GetBool("last") {
		log.Debugf("Since has been set to %s", rules.BeginningOfMonth(time.Now().AddDate(0, -1, 0)))
		since = rules.BeginningOfMonth(time.Now().AddDate(0, -1, 0))
		log.Debugf("Until has been set to %s", rules.EndOfMonth(time.Now().AddDate(0, -1, 0)))
		until = rules.EndOfMonth(time.Now().AddDate(0, -1, 0))
		return
	}

	log.Debugf("Since has been set to %s", rules.BeginningOfMonth(time.Now()))
	since = rules.BeginningOfMonth(time.Now())
	log.Debugf("Until has been set to %s", time.Now())
	until = time.Now()

	if sinceStr := viper.GetString("since"); sinceStr != "" {
		if sinceT, err := time.ParseInLocation("2006-01-02", sinceStr, config.GetLocation()); err != nil {
			log.WithError(err).Errorf("Failed to parse Since date %s", sinceStr)
		} else {
			log.Debugf("Since has been set to %s", sinceT)
			since = sinceT
		}
	}

	if untilStr := viper.GetString("until"); untilStr != "" {
		if untilT, err := time.ParseInLocation("2006-01-02", untilStr, config.GetLocation()); err != nil {
			log.WithError(err).Errorf("Failed to parse until date %s", untilStr)
		} else {
			log.Debugf("Until has been set to %s", untilT)
			until = untilT
		}
	}
}

func listOncallHours() {
	req := client.GET()
	req.Path("/oncalls")
	setSinceUntil()
	req.SetQuery("since", rules.FormatDate(since))
	req.SetQuery("until", rules.FormatDate(until))
	req.SetQuery("time_zone", config.GetTimeZone())
	if user := viper.GetString("user"); user != "" {
		if userID := client.GetUserID(user); userID != "" {
			req.SetQuery("user_ids[]", userID)
		}
	}

	if team := viper.GetString("team"); team != "" {
		if userIDs := client.GetAllUserID(team); userIDs != nil {
			for _, userID := range userIDs {
				req.AddQuery("user_ids[]", userID)
			}
		}
	}

	res, err := req.Do()
	if err != nil {
		log.WithError(err).Fatal("HTTP Client failed to GET Response from PagerDuty API")
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
			fmt.Printf("Start: %s - End: %s  Duration: %s\n", oncall.Start, oncall.End, oncall.End.Sub(oncall.Start))
			onWorkHours, offWorkHours := calculateWorkHours(oncall.Start, oncall.End, config.GetLocation())
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

func calculateWorkHours(start, end time.Time, location *time.Location) (int, int) {
	onWork, offWork := 0, 0
	elapsed := end.Sub(start)
	for offSet := time.Duration(0); offSet <= elapsed; offSet += time.Hour {
		// Check every hour (as a timestamp) from start to end
		sample := start.Add(offSet)
		// Skip hour that fell out of the since-until timeframe
		// This is needed because PagerDuty API returns the whole on-call intervalls from the queryed period
		// not just the parts that actually fell in the observed timefram.
		if sample.Before(since) || sample.After(until) {
			continue
		}
		if rules.IsWorkHour(sample, location) {
			onWork++
		} else {
			offWork++
		}
	}
	return onWork, offWork
}
