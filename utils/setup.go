package utils

import (
	"github.com/AngelVI13/slack-assistant/config"
	"github.com/AngelVI13/slack-assistant/slack/handlers"
	"github.com/AngelVI13/slack-assistant/users"
	"github.com/slack-go/slack/socketmode"
)

func SetupDeviceManager(config *config.Config, socketClient *socketmode.Client) *handlers.DeviceManager {
	devicesInfo := GetDevices(config)

	users := GetUsers(config.UsersFilename)
	deviceManager := &handlers.DeviceManager{
		DevicesMap:  devicesInfo,
		UsersInfo:   users,
		SlackClient: socketClient,
	}
	return deviceManager
}

func SetupReviewersList(config *config.Config, usersInfo users.UsersInfo) (reviewers []users.Reviewer) {
	for name, props := range usersInfo {
		if !props.IsReviewer {
			continue
		}

		reviewers = append(reviewers, users.Reviewer{Name: name, Id: props.Id})
	}
	return reviewers
}
