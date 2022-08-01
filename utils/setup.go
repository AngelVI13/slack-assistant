package utils

import (
	"github.com/AngelVI13/slack-assistant/config"
	"github.com/AngelVI13/slack-assistant/data"
	"github.com/AngelVI13/slack-assistant/users"
)

// SetupDataHolder Loads all data sources from locations specified in config file
func SetupDataHolder(config *config.Config) *data.DataHolder {
	devicesInfo := GetDevices(config)
	usersMap := GetUsers(config.UsersFilename)
	usersInfo := &users.UsersInfo{
		Map:      usersMap,
		Filename: config.UsersFilename,
	}
	parkingLot := GetParkingLot(config)

	dataHolder := &data.DataHolder{
		Devices:    &devicesInfo,
		Users:      usersInfo,
		Reviewers:  users.NewReviewers(config, &usersInfo.Map),
		ParkingLot: &parkingLot,
	}
	return dataHolder
}
