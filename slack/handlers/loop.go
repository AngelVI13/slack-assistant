package handlers

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/AngelVI13/slack-assistant/device"
	"github.com/AngelVI13/slack-assistant/slack/modals"
	"github.com/AngelVI13/slack-assistant/utils/users"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

<<<<<<< HEAD
var SlashCommands = map[string]ModalHandler{
	"/show-devices":   &modals.ShowDeviceHandler{},
	"/reserve-device": &modals.ReserveDeviceHandler{},
	"/release-device": &modals.ReleaseDeviceHandler{},
	"/restart-proxy":  &modals.RestartProxyHandler{},
	"/show-users":     &modals.ShowUsersHandler{},
=======
var SlashCommands = map[string]modals.ModalHandler{
	"/devices": modals.NewCustomOptionModalHandler(
		modals.DeviceActionMap,
		modals.DefaultDeviceAction,
		modals.DeviceModalInfo,
	),
	"/restart-proxy": &modals.RestartProxyHandler{},
>>>>>>> origin/master
}

type DeviceManager struct {
	device.DevicesMap
	// TODO: add possibility to extend this from slack
	Users       users.UserMap
	SlackClient *socketmode.Client
	// Whenever we are dealing with a modal that contains a state switching option
	// keep a pointer to it so we can change states
	CurrentOptionModalHandler modals.OptionModalHandler
}

func (dm *DeviceManager) GetDevicesInfo() device.DevicesInfo {
	devices := make(device.DevicesInfo, 0, len(dm.Devices))

	for _, value := range dm.Devices {
		devices = append(devices, value)
	}

	// NOTE: This sorts the device list starting from free devices
	sort.Slice(devices, func(i, j int) bool {
		return !devices[i].Reserved
	})

	firstTaken := -1 // Index of first taken device
	for i, device := range devices {
		if device.Reserved {
			firstTaken = i
			break
		}
	}

	// NOTE: this might be unnecessary but it shows devices in predicable way in UI so its nice.
	// If all devices are free or all devices are taken, sort by name
	if firstTaken == -1 || firstTaken == 0 {
		sort.Slice(devices, func(i, j int) bool {
			return devices[i].Name < devices[j].Name
		})
	} else {
		// split devices into 2 - free & taken
		// sort each sub slice based on device name/port
		free := devices[:firstTaken]
		taken := devices[firstTaken:]

		sort.Slice(free, func(i, j int) bool {
			return free[i].Name < free[j].Name
		})

		sort.Slice(taken, func(i, j int) bool {
			return taken[i].Name < taken[j].Name
		})
	}

	return devices
}

func (dm *DeviceManager) processEventApi(event socketmode.Event) {
	// The Event sent on the channel is not the same as the EventAPI events so we need to type cast it
	eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
	if !ok {
		log.Fatalf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
	}

	// We need to send an Acknowledge to the slack server
	dm.SlackClient.Ack(*event.Request)
	// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch
	err := HandleEventMessage(eventsAPIEvent, dm.SlackClient)
	if err != nil {
		// Replace with actual err handeling
		log.Fatal(err)
	}
}

func (dm *DeviceManager) processEventInteractive(event socketmode.Event) {
	interaction, ok := event.Data.(slack.InteractionCallback)
	if !ok {
		log.Fatalf("Could not type cast the message to a Interaction callback: %v\n", interaction)
	}

	err := dm.handleInteractionEvent(interaction)
	if err != nil {
		log.Fatal(err)
	}
	dm.SlackClient.Ack(*event.Request)
}

func (dm *DeviceManager) processSlashCommand(event socketmode.Event) {
	// Just like before, type cast to the correct event type, this time a SlashEvent
	command, ok := event.Data.(slack.SlashCommand)
	if !ok {
		log.Fatalf("Could not type cast the message to a SlashCommand: %v\n", command)
	}

	_, userAllowed := dm.Users[command.UserName]
	if !userAllowed {
		log.Printf("WARNING: Unauthorized user is sending command [%s] to the bot (%s)", command.Command, command.UserName)

		dm.handleUnauthorizedUserCommand(&command)
		dm.SlackClient.Ack(*event.Request, nil)
		return
	}

	log.Printf("PROCESS: Processing SLASH (%s) from (%s)", command.Command, command.UserName)
	// handleSlashCommand will take care of the command
	err := dm.handleSlashCommand(command)
	if err != nil {
		log.Fatal(err)
	}
	// Dont forget to acknowledge the request and send the payload
	// (we don't yet have any payload to send so we send nil)
	dm.SlackClient.Ack(*event.Request, nil)

}

func (dm *DeviceManager) ProcessMessageLoop(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)

	// Create a for loop that selects either the context cancellation or the events incomming
	for {
		select {
		// inscase context cancel is called exit the goroutine
		case <-ctx.Done():
			ticker.Stop()
			log.Println("Shutting down socketmode listener")
			return
		case event := <-dm.SlackClient.Events:
			// We have a new Events, let's type switch the event
			// Add more use cases here if you want to listen to other events.
			switch event.Type {
			case socketmode.EventTypeEventsAPI:
				// Handle mentions
				// NOTE: there is no user restriction for app mentions
				dm.processEventApi(event)
			case socketmode.EventTypeSlashCommand:
				dm.processSlashCommand(event)
			case socketmode.EventTypeInteractive:
				// Handle interaction events i.e. user voted in our poll etc.
				dm.processEventInteractive(event)
			}
		case <-ticker.C:
			// Auto release devices at midnight
			//                YYYY  M  D  H  M  S  NS timezone
			when := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
			dm.AutoRelease(when)
		}
	}
}

// handleUnauthorizedUserCommand will show an error message modal to user.
// Crashes in case the slack client could not open model view
func (dm *DeviceManager) handleUnauthorizedUserCommand(command *slack.SlashCommand) {
	handler := &modals.UnauthorizedHandler{}
	modalRequest := handler.GenerateModalRequest(command, dm.GetDevicesInfo())

	_, err := dm.SlackClient.OpenView(command.TriggerID, modalRequest)
	if err != nil {
		log.Fatalf("Error opening view: %s", err)
	}
}

// handleSlashCommand will take a slash command and route to the appropriate function
func (dm *DeviceManager) handleSlashCommand(command slack.SlashCommand) error {
	handler, hasValue := SlashCommands[command.Command]
	if !hasValue {
		// NOTE: this can only happen if slack added new command but the bot was not updated to support it
		log.Printf("WARNING: User (%s) requested unsupported command %s\n", command.UserName, command.Command)
		return nil
	}
	return dm.handleDeviceCommand(&command, handler)
}

func (dm *DeviceManager) handleDeviceCommand(
	command *slack.SlashCommand,
	handler modals.ModalHandler,
) error {
<<<<<<< HEAD
	var data any
	if command.Command == "/show-users" {
		data = dm.Users
	} else {
		data = dm.GetDevicesInfo()
	}
	modalRequest := handler.GenerateModalRequest(command, data)
=======
	// In case we are dealing with an OptionModalHandler save pointer to it
	// so we can change its state when needed
	optionHandler, ok := handler.(modals.OptionModalHandler)
	if ok {
		dm.CurrentOptionModalHandler = optionHandler
		dm.CurrentOptionModalHandler.Reset()
	} else {
		dm.CurrentOptionModalHandler = nil
	}

	modalRequest := handler.GenerateModalRequest(dm.GetDevicesInfo())
>>>>>>> origin/master
	_, err := dm.SlackClient.OpenView(command.TriggerID, modalRequest)
	if err != nil {
		return fmt.Errorf("Error opening view: %s", err)
	}
	return nil
}

func (dm *DeviceManager) handleInteractionEvent(interaction slack.InteractionCallback) error {
	log.Println(interaction)
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
			cmdOutput := dm.RestartProxies(deviceNames, interaction.User.Name)
			dm.SlackClient.PostEphemeral(interaction.User.ID, interaction.User.ID, slack.MsgOptionText(cmdOutput, false))
		default:
		}

	case slack.InteractionTypeBlockActions:
		switch interaction.View.Title.Text {
		case modals.MDeviceTitle:
			if dm.CurrentOptionModalHandler == nil {
				log.Fatalf(
					`Did not have a valid pointer to OptionModal,
        				please make sure to close any open modals before restarting the bot`,
				)
			}

			// Update option view if new option was chosen
			option := interaction.View.State.Values[modals.MDeviceActionId][modals.MDeviceOptionId].SelectedOption.Value
			dm.CurrentOptionModalHandler.ChangeAction(option)

			// handle button actions
			for _, action := range interaction.ActionCallback.BlockActions {
				switch action.ActionID {
				case modals.ReserveDeviceActionId, modals.ReserveWithAutoActionId:

					autoRelease := action.ActionID == modals.ReserveWithAutoActionId
					errStr := dm.Reserve(action.Value, interaction.User.Name, interaction.User.ID, autoRelease)
					if errStr != "" {
						log.Println(errStr)
						// If there device was already taken -> inform user by personal DM message from the bot
						dm.SlackClient.PostEphemeral(interaction.User.ID, interaction.User.ID, slack.MsgOptionText(errStr, false))
					}
				case modals.ReleaseDeviceActionId:
					victimId, errStr := dm.Release(action.Value, interaction.User.Name)
					if victimId != "" {
						log.Println(errStr)
						dm.SlackClient.PostEphemeral(victimId, victimId, slack.MsgOptionText(errStr, false))
					}
				default:
				}
			}

			// update modal view to display changes
			updatedView := dm.CurrentOptionModalHandler.GenerateModalRequest(dm.GetDevicesInfo())
			_, err := dm.SlackClient.UpdateView(updatedView, "", "", interaction.View.ID)
			if err != nil {
				log.Fatal(err)
			}
		}
	default:

	}

	return nil
}
