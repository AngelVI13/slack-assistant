package utils

import (
	"github.com/AngelVI13/slack-assistant/config"
	"github.com/AngelVI13/slack-assistant/slack/handlers"
	"github.com/AngelVI13/slack-assistant/users"
)

// SetupDataHolder Loads all data sources from locations specified in config file
func SetupDataHolder(config *config.Config) *handlers.DataHolder {
	devicesInfo := GetDevices(config)
	usersMap := GetUsers(config.UsersFilename)
	usersInfo := &users.UsersInfo{
		Map:      usersMap,
		Filename: config.UsersFilename,
	}

	dataHolder := &handlers.DataHolder{
		Devices:   &devicesInfo,
		Users:     usersInfo,
		Reviewers: users.NewReviewers(config, &usersInfo.Map),
	}
	return dataHolder
}
