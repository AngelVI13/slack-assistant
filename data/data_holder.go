package data

import (
	"github.com/AngelVI13/slack-assistant/device"
	"github.com/AngelVI13/slack-assistant/users"
)

type DataHolder struct {
	Devices   *device.DevicesMap
	Users     *users.UsersInfo
	Reviewers users.Reviewers
}
