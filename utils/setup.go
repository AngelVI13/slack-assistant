package utils

import (
	"github.com/AngelVI13/slack-assistant/config"
	"github.com/AngelVI13/slack-assistant/slack/handlers"
	"github.com/AngelVI13/slack-assistant/users"
)

// SetupDataHolder Loads all data sources from locations specified in config file
func SetupDataHolder(config *config.Config) *handlers.DataHolder {
	devicesInfo := GetDevices(config)
	usersInfo := GetUsers(config.UsersFilename)
	reviewers := SetupReviewersList(usersInfo)

	dataHolder := &handlers.DataHolder{
		Devices:   &devicesInfo,
		Users:     usersInfo,
		Reviewers: reviewers,
	}
	return dataHolder
}

func SetupReviewersList(usersInfo users.UsersInfo) (reviewers []*users.Reviewer) {
	for name, props := range usersInfo {
		if !props.IsReviewer {
			continue
		}

		reviewers = append(reviewers, &users.Reviewer{Name: name, Id: props.Id})
	}
	return reviewers
}
