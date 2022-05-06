package config

import (
	"fmt"
	"os"
)

type Config struct {
	SlackAuthToken string
	SlackChannelId string
	SlackAppToken  string

	DevicesFilename string
	UsersFilename   string
	Debug           bool
	TaEndpoint      string
	WorkersEndpoint string
	ProxyEndpoint   string
	ProxyUrl        string
}

// ConfigFromEnv Creates config instance by reading corresponding ENV variables.
// Make sure godotenv.Load is called beforehand
func ConfigFromEnv() *Config {
	taEndpoint := os.Getenv("SL_TA_ENDPOINT")

	return &Config{
		SlackAuthToken: os.Getenv("SLACK_AUTH_TOKEN"),
		SlackChannelId: os.Getenv("SLACK_CHANNEL_ID"),
		SlackAppToken:  os.Getenv("SLACK_APP_TOKEN"),

		DevicesFilename: os.Getenv("SL_DEVICES_FILE"),
		UsersFilename:   os.Getenv("SL_USERS_FILE"),
		Debug:           os.Getenv("SL_DEBUG") == "1",
		TaEndpoint:      taEndpoint,
		WorkersEndpoint: fmt.Sprintf("%s/workers", taEndpoint),
		ProxyEndpoint:   fmt.Sprintf("%s/proxy", taEndpoint),
		ProxyUrl:        os.Getenv("SL_PROXY"),
	}
}
