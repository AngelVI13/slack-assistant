package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AngelVI13/slack-assistant/device"
	"github.com/AngelVI13/slack-assistant/parking"
	"github.com/AngelVI13/slack-assistant/slack/modals"
	"github.com/AngelVI13/slack-assistant/slack/slash"
	"github.com/AngelVI13/slack-assistant/users"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

var SlashCommandsForModals = map[string]modals.ModalHandler{
	"/devices": modals.NewCustomOptionModalHandler(
		modals.DeviceActionMap,
		modals.DefaultDeviceAction,
		modals.DeviceModalInfo,
	),
	"/restart-proxy": &modals.RestartProxyHandler{},
	"/users": modals.NewCustomOptionModalHandler(
		modals.UsersActionMap,
		modals.DefaultUsersAction,
		modals.UsersModalInfo,
	),
	"/test-park": modals.NewCustomOptionModalHandler(
		modals.ParkingActionMap,
		modals.DefaultParkingAction,
		modals.ParkingModalInfo,
	),
}

var SlashCommandsForHandlers = map[string]slash.SlashHandler{
	"/review": &slash.ReviewHandler{},
}

// TODO: move this somewhere else and not in loop file
type DataHolder struct {
	Devices    *device.DevicesMap
	Users      *users.UsersInfo
	Reviewers  users.Reviewers
	ParkingLot *parking.ParkingLot
}

type SlackBot struct {
	Data *DataHolder

	SlackClient *socketmode.Client
	// Whenever we are dealing with a modal that contains a state switching option
	// keep a pointer to it so we can change states
	CurrentOptionModalHandler modals.OptionModalHandler
}

func (bot *SlackBot) processEventApi(event socketmode.Event) {
	// The Event sent on the channel is not the same as the EventAPI events so we need to type cast it
	eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
	if !ok {
		log.Fatalf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
	}

	// We need to send an Acknowledge to the slack server
	bot.SlackClient.Ack(*event.Request)
	// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch
	err := HandleEventMessage(eventsAPIEvent, bot.SlackClient)
	if err != nil {
		// Replace with actual err handeling
		log.Fatal(err)
	}
}

func (bot *SlackBot) processEventInteractive(event socketmode.Event) {
	interaction, ok := event.Data.(slack.InteractionCallback)
	if !ok {
		log.Fatalf("Could not type cast the message to a Interaction callback: %v\n", interaction)
	}

	err := bot.handleInteractionEvent(interaction)
	if err != nil {
		log.Fatal(err)
	}
	bot.SlackClient.Ack(*event.Request)
}

func (bot *SlackBot) processSlashCommand(event socketmode.Event) {
	// Just like before, type cast to the correct event type, this time a SlashEvent
	command, ok := event.Data.(slack.SlashCommand)
	if !ok {
		log.Fatalf("Could not type cast the message to a SlashCommand: %v\n", command)
	}

	_, userAllowed := bot.Data.Users.Map[command.UserName]
	if !userAllowed {
		log.Printf("WARNING: Unauthorized user is sending command [%s] to the bot (%s)", command.Command, command.UserName)

		bot.handleUnauthorizedUserCommand(&command)
		bot.SlackClient.Ack(*event.Request, nil)
		return
	}

	log.Printf("PROCESS: Processing SLASH (%s) from (%s)", command.Command, command.UserName)
	// handleSlashCommand will take care of the command
	err := bot.handleSlashCommand(command)
	if err != nil {
		log.Fatal(err)
	}
	// Dont forget to acknowledge the request and send the payload
	// (we don't yet have any payload to send so we send nil)
	bot.SlackClient.Ack(*event.Request, nil)

}

func (bot *SlackBot) ProcessMessageLoop(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)

	// Create a for loop that selects either the context cancellation or the events incomming
	for {
		select {
		// inscase context cancel is called exit the goroutine
		case <-ctx.Done():
			ticker.Stop()
			log.Println("Shutting down socketmode listener")
			return
		case event := <-bot.SlackClient.Events:
			// We have a new Events, let's type switch the event
			// Add more use cases here if you want to listen to other events.
			switch event.Type {
			case socketmode.EventTypeEventsAPI:
				// Handle mentions
				// NOTE: there is no user restriction for app mentions
				bot.processEventApi(event)
			case socketmode.EventTypeSlashCommand:
				bot.processSlashCommand(event)
			case socketmode.EventTypeInteractive:
				// Handle interaction events i.e. user voted in our poll etc.
				bot.processEventInteractive(event)
			}
		case <-ticker.C:
			// Auto release devices at midnight
			//                YYYY  M  D  H  M  S  NS timezone
			when := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
			bot.Data.Devices.AutoRelease(when)
		}
	}
}

// handleUnauthorizedUserCommand will show an error message modal to user.
// Crashes in case the slack client could not open model view
func (bot *SlackBot) handleUnauthorizedUserCommand(command *slack.SlashCommand) {
	handler := &modals.UnauthorizedHandler{}
	modalRequest := handler.GenerateModalRequest(command, bot.Data.Devices.GetDevicesInfo())

	_, err := bot.SlackClient.OpenView(command.TriggerID, modalRequest)
	if err != nil {
		log.Fatalf("Error opening view: %s", err)
	}
}

// handleSlashCommand will take a slash command and route to the appropriate function
func (bot *SlackBot) handleSlashCommand(command slack.SlashCommand) error {
	if handler, hasValue := SlashCommandsForModals[command.Command]; hasValue {
		return bot.handleDeviceCommand(&command, handler)
	} else if handler, hasValue := SlashCommandsForHandlers[command.Command]; hasValue {
		// TODO: Reviewers here is hardcoded -> need a better way to handle args for slash commands
		return handler.Execute(&command, bot.SlackClient, &bot.Data.Reviewers)
	} else {
		// NOTE: this can only happen if slack added new command but the bot was not updated to support it
		log.Printf("WARNING: User (%s) requested unsupported command %s\n", command.UserName, command.Command)
		return nil
	}

}

func (bot *SlackBot) handleDeviceCommand(
	command *slack.SlashCommand,
	handler modals.ModalHandler,
) error {

	// TODO: fix this
	var data any
	if command.Command == "/test-users" {
		data = bot.Data.Users.Map
	} else if command.Command == "/test-park" {
		data = bot.Data.ParkingLot.GetSpacesInfo()
	} else {
		data = bot.Data.Devices.GetDevicesInfo()
	}

	// In case we are dealing with an OptionModalHandler save pointer to it
	// so we can change its state when needed
	optionHandler, ok := handler.(modals.OptionModalHandler)
	if ok {
		bot.CurrentOptionModalHandler = optionHandler
		bot.CurrentOptionModalHandler.Reset()
	} else {
		bot.CurrentOptionModalHandler = nil
	}

	modalRequest := handler.GenerateModalRequest(data)
	_, err := bot.SlackClient.OpenView(command.TriggerID, modalRequest)
	if err != nil {
		return fmt.Errorf("Error opening view: %s", err)
	}
	return nil
}

func (bot *SlackBot) handleInteractionEvent(interaction slack.InteractionCallback) error {
	switch interaction.Type {
	case slack.InteractionTypeViewSubmission:
		// NOTE: we use title text to determine which modal was submitted
		switch interaction.View.Title.Text {
		case modals.MRestartProxyTitle:
			deviceNames := []string{}
			userSelection := interaction.View.State.Values[modals.MRestartProxyActionId][modals.MRestartProxyCheckboxId].SelectedOptions
			for _, selected := range userSelection {
				deviceNames = append(deviceNames, selected.Value)
			}
			cmdOutput := bot.Data.Devices.RestartProxies(deviceNames, interaction.User.Name)
			bot.SlackClient.PostEphemeral(interaction.User.ID, interaction.User.ID, slack.MsgOptionText(cmdOutput, false))
		case modals.MRemoveUsersTitle:
			userSelection := interaction.View.State.Values[modals.MRemoveUsersActionId][modals.MRemoveUsersOptionId].SelectedOptions

			for name := range bot.Data.Users.Map {
				for _, a := range userSelection {
					if a.Value == name {
						log.Printf("Deleting %s", a.Value)
						delete(bot.Data.Users.Map, a.Value)
					}
				}
			}

			bot.Data.Users.SynchronizeToFile()

		case modals.MAddUserTitle:

			userSelection := interaction.View.State.Values[modals.MAddUserActionId][modals.MAddUserOptionId].SelectedUsers

			for _, new_user := range userSelection {
				user_info, _ := bot.SlackClient.GetUserInfo(new_user)
				user_name := user_info.Name
				log.Printf("Adding %s", user_name)
				// TODO: get access rights and isReviewer from input
				bot.Data.Users.Map[user_name] = &users.User{
					Id:         user_info.ID,
					Rights:     users.STANDARD, // Assign only standart value for now
					IsReviewer: false,
				}
			}

			bot.Data.Users.SynchronizeToFile()
		default:
		}

	case slack.InteractionTypeBlockActions:
		switch interaction.View.Title.Text {
		case modals.MDeviceTitle:
			if bot.CurrentOptionModalHandler == nil {
				log.Fatalf(
					`Did not have a valid pointer to OptionModal,
        				please make sure to close any open modals before restarting the bot`,
				)
			}

			// Update option view if new option was chosen
			option := interaction.View.State.Values[modals.MDeviceActionId][modals.MDeviceOptionId].SelectedOption.Value
			bot.CurrentOptionModalHandler.ChangeAction(option)

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
			updatedView := bot.CurrentOptionModalHandler.GenerateModalRequest(bot.Data.Devices.GetDevicesInfo())
			_, err := bot.SlackClient.UpdateView(updatedView, "", "", interaction.View.ID)
			if err != nil {
				log.Fatal(err)
			}
		case modals.MShowUsersTitle, modals.MRemoveUsersTitle, modals.MAddUserTitle:
			if bot.CurrentOptionModalHandler == nil {
				log.Fatalf(
					`Did not have a valid pointer to OptionModal,
						please make sure to close any open modals before restarting the bot`,
				)
			}

			// Update option view if new option was chosen
			option := interaction.View.State.Values[modals.MUsersActionId][modals.MUsersOptionId].SelectedOption.Value
			log.Println(option)
			ok := bot.CurrentOptionModalHandler.ChangeAction(option)
			log.Println(ok)

			// update modal view to display changes
			updatedView := bot.CurrentOptionModalHandler.GenerateModalRequest(bot.Data.Users.Map)
			_, err := bot.SlackClient.UpdateView(updatedView, "", "", interaction.View.ID)
			if err != nil {
				log.Fatal(err)
			}

		case modals.MParkingTitle:
			if bot.CurrentOptionModalHandler == nil {
				log.Fatalf(
					`Did not have a valid pointer to OptionModal,
        				please make sure to close any open modals before restarting the bot`,
				)
			}

			// Update option view if new option was chosen
			option := interaction.View.State.Values[modals.MParkingActionId][modals.MParkingOptionId].SelectedOption.Value
			bot.CurrentOptionModalHandler.ChangeAction(option)

			// handle button actions
			for _, action := range interaction.ActionCallback.BlockActions {
				switch action.ActionID {
				case modals.ReserveParkingActionId, modals.ReserveParkingWithAutoActionId:
					autoRelease := action.ActionID == modals.ReserveParkingWithAutoActionId
					errStr := bot.Data.ParkingLot.Reserve(action.Value, interaction.User.Name, interaction.User.ID, autoRelease)
					if errStr != "" {
						log.Println(errStr)
						// If there device was already taken -> inform user by personal DM message from the bot
						bot.SlackClient.PostEphemeral(interaction.User.ID, interaction.User.ID, slack.MsgOptionText(errStr, false))
					}
				case modals.ReleaseParkingActionId:
					victimId, errStr := bot.Data.ParkingLot.Release(action.Value, interaction.User.Name)
					if victimId != "" {
						log.Println(errStr)
						bot.SlackClient.PostEphemeral(victimId, victimId, slack.MsgOptionText(errStr, false))
					}
				default:
				}
			}

			// update modal view to display changes
			updatedView := bot.CurrentOptionModalHandler.GenerateModalRequest(bot.Data.ParkingLot.GetSpacesInfo())
			_, err := bot.SlackClient.UpdateView(updatedView, "", "", interaction.View.ID)
			if err != nil {
				log.Fatal(err)
			}
		default:
		}
	default:

	}

	return nil
}
