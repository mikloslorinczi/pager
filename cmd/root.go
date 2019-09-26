package cmd

import (
	"github.com/mikloslorinczi/pager/config"
	"github.com/mikloslorinczi/pager/logger"
	"github.com/spf13/cobra"
)

var (
	log = logger.NewLogger("Cmd")
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Long: `
Pager is a command line helper for Pager Duty Schedules, Oncalls and more...
`,
	PersistentPreRun: config.Setup,
}

func init() {
	// --config is the only flag which is not bind to a Viper key
	// this is because the config file is needed during the setup of Viper
	rootCmd.PersistentFlags().StringVarP(&config.ConfigFile, "config", "c", "", "Config file")
	rootCmd.PersistentFlags().StringP("api_url", "u", "https://api.pagerduty.com", "Pager Duty API URL")
	rootCmd.PersistentFlags().StringP("api_token", "t", "", "Pager Duty API Token")
}

// Execute is the entrypoint of the Cobra app, it executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.ErrorAndStacktrace(err, "Error during execution")
	}
}
