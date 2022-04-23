package utils

import (
	"github.com/AngelVI13/slack-assistant/config"
	"github.com/AngelVI13/slack-assistant/slack/handlers"
	"github.com/slack-go/slack/socketmode"
)

func SetupDeviceManager(config *config.Config, socketClient *socketmode.Client) *handlers.DeviceManager {
	devicesInfo := GetDevices(config)

	users := GetUsers(config.UsersFilename)
	deviceManager := &handlers.DeviceManager{
		DevicesMap:  devicesInfo,
		Users:       users,
		SlackClient: socketClient,
	}
	return deviceManager
}
