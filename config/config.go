package config

import (
	"github.com/mikloslorinczi/pager/client"
	"github.com/mikloslorinczi/pager/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// ConfigFile is the path to the Testman config file it can be set with the --path (-p) flag
	ConfigFile string
	log        = logger.NewLogger("Config")
)

// Setup sets up Viper and the HTTP Client
func Setup(cmd *cobra.Command, args []string) {
	SetViper(cmd, args)
	log.Debug("Viper has been set up")

	if len(viper.GetString("api_token")) < 1 {
		log.Error("No API Token specified! You can set it with the --api_token flag or the PAGER_API_TOKEN env var")
	}
	client.Setup()
	log.Debug("HTTP Client has been set up")
}

// SetViper should be the PersistentPreRun of all Cobra command (attached at the root command level)
// It considers defaults, environment variables prefixed with PAGER_ and Cobra flags
// (in this order from weaker to stronger) and bind them to Viper keys
func SetViper(cmd *cobra.Command, args []string) {
	// Bind Cobra flags to Viper keys
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		log.WithError(err).Fatal("Failed to bind Cobra flags to Viper keys")
	}
	// Search the environment for variables prefixed with PAGER_
	viper.SetEnvPrefix("PAGER")
	// Read in environment variables that match
	viper.AutomaticEnv()
	// If a ConfigFile is set read it in
	if ConfigFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(ConfigFile)
		if err := viper.ReadInConfig(); err != nil {
			log.WithError(err).Fatalf("Failed to load Viper configs from file %s", ConfigFile)
		}
	}
}
