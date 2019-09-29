package cmd

import (
	"fmt"

	"github.com/mikloslorinczi/pager/client"
	"github.com/mikloslorinczi/pager/config"
	"github.com/mikloslorinczi/pager/model"

	"github.com/spf13/cobra"
)

// oncallCmd represents the run command
var oncallCmd = &cobra.Command{
	Use:     "oncall",
	Aliases: []string{"oc", "onc", "call", "on-call", "oncalls", "on-calls"},
	Short:   "On-call information",
	Long: `
Oncall information
Who and when were are and will be on call
will be told by this shiny christal ball
	`,
	Example: "pager -t PD_TOKEN oncall list",
	Run: func(cmd *cobra.Command, args []string) {
		onCall()
	},
}

func init() {
	rootCmd.AddCommand(oncallCmd)
}

func onCall() {
	req := client.GET()
	req.Path("/oncalls")
	req.AddQuery("time_zone", config.GetTimeZone())
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
	for _, oncall := range oncallsResp.Oncalls {
		fmt.Printf("User: %s\n", oncall.User.Summary)
		fmt.Printf("Escalation Policy: %s\n", oncall.EscalationPolicy.Summary)
		fmt.Printf("Escalation Level: %d\n", oncall.EscalationLevel)
		if schedule := oncall.Schedule.Summary; len(schedule) > 0 {
			fmt.Printf("Schedule: %s\n", oncall.Schedule.Summary)
			fmt.Printf("Start: %s - End: %s\n", oncall.Start, oncall.End)
		}
		fmt.Println("---------------------------------------------------------")
	}
}
