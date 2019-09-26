package client

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"

	"gopkg.in/h2non/gentleman.v2"
)

var (
	defaultClient *gentleman.Client
	once          sync.Once
)

// Setup sets up the defaultClient with the Pager Duty URL and Token provided by Viper
func Setup() {
	once.Do(func() {
		defaultClient = gentleman.New()
		defaultClient.URL(viper.GetString("api_url"))
		defaultClient.SetHeader("Accept", "application/json")
		defaultClient.SetHeader("Authorization", fmt.Sprintf("Token token=%s", viper.GetString("api_token")))
	})
}

// GET creates a new GET Request with the defaultClient and returns it
func GET() *gentleman.Request {
	return defaultClient.Get()
}
