package config

import (
	"fmt"
	"os"
)

type Config struct {
	SlackAuthToken string
	SlackChannelId string
	SlackAppToken  string

	DevicesFilename   string
	UsersFilename     string
	ReviewersFilename string

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

		// NOTE: this file is used to store current list of reviewers
		// i.e. reviewers are selected one by one until everyone has taken his turn
		// after which the list is reset to full reviewers list.
		// I don't see a reason why you might want to have that filename configurable
		// so hardcoded it will stay.
		ReviewersFilename: ".reviewers.txt",

		Debug:           os.Getenv("SL_DEBUG") == "1",
		TaEndpoint:      taEndpoint,
		WorkersEndpoint: fmt.Sprintf("%s/workers", taEndpoint),
		ProxyEndpoint:   fmt.Sprintf("%s/proxy", taEndpoint),
		ProxyUrl:        os.Getenv("SL_PROXY"),
	}
}
