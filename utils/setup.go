package utils

import (
	"fmt"
	"os"

	"github.com/AngelVI13/slack-assistant/slack/handlers"
	"github.com/slack-go/slack/socketmode"
)

func SetupDeviceManager(socketClient *socketmode.Client) *handlers.DeviceManager {
	devicesFile := os.Getenv("SL_DEVICES_FILE")
	workersEndpoint := fmt.Sprintf("%s/workers", os.Getenv("SL_TA_ENDPOINT"))
	devicesInfo := GetDevices(devicesFile, workersEndpoint)

	usersFile := os.Getenv("SL_USERS_FILE")
	users := GetUsers(usersFile)
	deviceManager := &handlers.DeviceManager{
		DevicesMap:  devicesInfo,
		Users:       users,
		SlackClient: socketClient,
	}
	return deviceManager
}
