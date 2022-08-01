package handlers

import (
	"log"

	"github.com/AngelVI13/slack-assistant/device"
	"github.com/AngelVI13/slack-assistant/slack/modals"
	"github.com/AngelVI13/slack-assistant/users"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func addUserSubmission(
	slackClient *socketmode.Client,
	interaction *slack.InteractionCallback,
	usersInfo *users.UsersInfo,
	reviewers *users.Reviewers,
) {
	selectedUsers := interaction.View.State.Values[modals.MAddUserActionId][modals.MAddUserOptionId].SelectedUsers
	selectedOptions := interaction.View.State.Values[modals.MAddUserAccessRightActionId][modals.MAddUserAccessRightOptionId].SelectedOptions

	selectedUsersInfo := []*slack.User{}
	for _, newUser := range selectedUsers {
		// TODO: this seems wrong but i think someone else is working on it and its not finished
		userInfo, _ := slackClient.GetUserInfo(newUser)
		selectedUsersInfo = append(selectedUsersInfo, userInfo)
	}

	usersInfo.AddNewUsers(selectedUsersInfo, selectedOptions, modals.MAddUserAccessRightOption, modals.MAddUserReviewerOption)
	reviewers.All = users.GetReviewers(&usersInfo.Map)
}

func removeUserSubmission(interaction *slack.InteractionCallback, usersInfo *users.UsersInfo) {
	userSelection := interaction.View.State.Values[modals.MRemoveUsersActionId][modals.MRemoveUsersOptionId].SelectedOptions

	for name := range usersInfo.Map {
		for _, a := range userSelection {
			if a.Value == name {
				log.Printf("Deleting %s", a.Value)
				delete(usersInfo.Map, a.Value)
			}
		}
	}

	usersInfo.SynchronizeToFile()
}

func restartProxySubmission(slackClient *socketmode.Client, interaction *slack.InteractionCallback, devices *device.DevicesMap) {
	deviceNames := []string{}
	userSelection := interaction.View.State.Values[modals.MRestartProxyActionId][modals.MRestartProxyCheckboxId].SelectedOptions
	for _, selected := range userSelection {
		deviceNames = append(deviceNames, selected.Value)
	}
	cmdOutput := devices.RestartProxies(deviceNames, interaction.User.Name)
	slackClient.PostEphemeral(interaction.User.ID, interaction.User.ID, slack.MsgOptionText(cmdOutput, false))
}
