package config

import (
	"time"

	"github.com/mikloslorinczi/pager/client"
	"github.com/mikloslorinczi/pager/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// ConfigFile is the path to the Testman config file it can be set with the --path (-p) flag
	ConfigFile string
	timeZone   *time.Location
	log        = logger.NewLogger("Config")
)

// Setup sets up Viper, the HTTP Client and the TimeZone
func Setup(cmd *cobra.Command, args []string) {
	setViper(cmd, args)
	log.Debug("Viper has been set up")

	if len(viper.GetString("api_token")) < 1 {
		log.Error("No API Token specified! You can set it with the --api_token flag or the PAGER_API_TOKEN env var")
	}
	client.Setup()
	log.Debug("HTTP Client has been set up")

	setTimeZone()
	log.Debugf("Time Zone has been set to %s", GetTimeZone())
}

// setViper should be the PersistentPreRun of all Cobra command (attached at the root command level)
// It considers defaults, environment variables prefixed with PAGER_ and Cobra flags
// (in this order from weaker to stronger) and bind them to Viper keys
func setViper(cmd *cobra.Command, args []string) {
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

// setTimeZone sets the Time Zone according to Viper (if it is not set in Viper defaults to CET)
func setTimeZone() {
	if zone := viper.GetString("time_zone"); zone != "" {
		if loc, err := time.LoadLocation(zone); err != nil {
			log.WithError(err).Errorf("Failed to parse Time Zone %s", zone)
		} else {
			timeZone = loc
			return
		}
	}
	timeZone = time.Local
}

// GetTimeZone returns the string representation of the Time Zone
func GetTimeZone() string {
	return timeZone.String()
}

// GetLocation returns the Time Zone location
func GetLocation() *time.Location {
	return timeZone
}
