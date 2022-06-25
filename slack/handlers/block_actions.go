package handlers

import (
	"log"

	"github.com/AngelVI13/slack-assistant/slack/modals"
	"github.com/slack-go/slack"
)

func handleDeviceActions(bot *SlackBot, interaction *slack.InteractionCallback) *slack.ModalViewRequest {
	// Update option view if new option was chosen
	option := interaction.View.State.Values[modals.MDeviceActionId][modals.MDeviceOptionId].SelectedOption.Value
	bot.CurrentOptionModalData.Handler.ChangeAction(option)

	// handle button actions
	for _, action := range interaction.ActionCallback.BlockActions {
		switch action.ActionID {
		case modals.ReserveDeviceActionId, modals.ReserveWithAutoActionId:

			autoRelease := action.ActionID == modals.ReserveWithAutoActionId
			errStr := bot.Data.Devices.Reserve(action.Value, interaction.User.Name, interaction.User.ID, autoRelease)
			if errStr != "" {
				log.Println(errStr)
				// If there device was already taken -> inform user by personal DM message from the bot
				bot.SlackClient.PostEphemeral(interaction.User.ID, interaction.User.ID, slack.MsgOptionText(errStr, false))
			}
		case modals.ReleaseDeviceActionId:
			victimId, errStr := bot.Data.Devices.Release(action.Value, interaction.User.Name)
			if victimId != "" {
				log.Println(errStr)
				bot.SlackClient.PostEphemeral(victimId, victimId, slack.MsgOptionText(errStr, false))
			}
		default:
		}
	}

	// update modal view to display changes
	updatedView := bot.CurrentOptionModalData.Handler.GenerateModalRequest(
		bot.CurrentOptionModalData.Command,
		bot.Data.Devices.GetDevicesInfo(
			bot.CurrentOptionModalData.Command.UserName,
		),
	)
	return &updatedView
}

func handleUserActions(bot *SlackBot, interaction *slack.InteractionCallback) *slack.ModalViewRequest {
	// Update option view if new option was chosen
	option := interaction.View.State.Values[modals.MUsersActionId][modals.MUsersOptionId].SelectedOption.Value
	log.Println(option)
	ok := bot.CurrentOptionModalData.Handler.ChangeAction(option)
	log.Println(ok)

	updatedView := bot.CurrentOptionModalData.Handler.GenerateModalRequest(bot.CurrentOptionModalData.Command, bot.Data.Users.Map)
	return &updatedView
}

func handleParkingActions(bot *SlackBot, interaction *slack.InteractionCallback) *slack.ModalViewRequest {
	// Update option view if new option was chosen
	option := interaction.View.State.Values[modals.MParkingActionId][modals.MParkingOptionId].SelectedOption.Value
	bot.CurrentOptionModalData.Handler.ChangeAction(option)

	// Check if an admin has made the request
	isSpecialUser := bot.Data.Users.IsSpecial(interaction.User.Name)

	// handle button actions
	for _, action := range interaction.ActionCallback.BlockActions {
		parkingSpace := action.Value
		switch action.ActionID {
		case modals.ReserveParkingActionId:
			handleReserveParking(bot, interaction, parkingSpace, isSpecialUser)
		case modals.ReleaseParkingActionId:
			handleReleaseParking(bot, interaction, parkingSpace, isSpecialUser)
		default:
		}
	}

	// update modal view to display changes
	updatedView := bot.CurrentOptionModalData.Handler.GenerateModalRequest(
		bot.CurrentOptionModalData.Command,
		bot.Data.ParkingLot.GetSpacesInfo(bot.CurrentOptionModalData.Command.UserName),
	)
	return &updatedView
}

func handleReserveParking(bot *SlackBot, interaction *slack.InteractionCallback, parkingSpace string, isSpecialUser bool) {
	autoRelease := true // by default parking reservation is always with auto release
	if isSpecialUser {  // unless we have a special user (i.e. user with designated parking space)
		autoRelease = false
	}

	errStr := bot.Data.ParkingLot.Reserve(parkingSpace, interaction.User.Name, interaction.User.ID, autoRelease)
	if errStr != "" {
		log.Println(errStr)
		// If there device was already taken -> inform user by personal DM message from the bot
		bot.SlackClient.PostEphemeral(interaction.User.ID, interaction.User.ID, slack.MsgOptionText(errStr, false))
	}
}

func handleReleaseParking(bot *SlackBot, interaction *slack.InteractionCallback, parkingSpace string, isSpecialUser bool) {
	if !isSpecialUser {
		victimId, errStr := bot.Data.ParkingLot.Release(parkingSpace, interaction.User.Name)
		if victimId != "" {
			log.Println(errStr)
			bot.SlackClient.PostEphemeral(victimId, victimId, slack.MsgOptionText(errStr, false))
		}
	} else {
		chosenParkingSpace := bot.Data.ParkingLot.GetSpace(parkingSpace)

		parkingReleaseHandler := &modals.ParkingReleaseHandler{}
		updatedView := parkingReleaseHandler.GenerateModalRequest(
			bot.CurrentOptionModalData.Command,
			chosenParkingSpace,
		)
		log.Println("_____ Generating BOOKING view")
		_, err := bot.SlackClient.PushView(interaction.TriggerID, updatedView)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func handleParkingBooking(bot *SlackBot, interaction *slack.InteractionCallback) *slack.ModalViewRequest {
	// handle button actions
	for _, action := range interaction.ActionCallback.BlockActions {
		switch action.ActionID {
		case modals.ReleaseStartDateActionId:
			// format is YYYY-MM-DD
			log.Println("----------- Start Date ", action.SelectedDate)
		case modals.ReleaseEndDateActionId:
			log.Println("----------- End Date ", action.SelectedDate)
		default:
		}
	}

	/*
		parkingReleaseHandler := &modals.ParkingReleaseHandler{}
		updatedView := parkingReleaseHandler.GenerateModalRequest(
			bot.CurrentOptionModalData.Command,
			chosenParkingSpace,
		)
		_, err := bot.SlackClient.UpdateView(updatedView, "", "", interaction.View.ID)
		if err != nil {
			log.Fatal(err)
		}
	*/
	return nil
}
