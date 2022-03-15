package handlers

import (
	"context"
	"fmt"
	"github.com/AngelVI13/slack-assistant/modals"
	"log"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

const SlashShowDevices = "/show-devices"
const SlashReserveDevice = "/reserve-device"
const SlashReleaseDevice = "/release-device"

var SlashCommands = map[string]ModalHandler{
	SlashShowDevices:   &modals.ShowDeviceHandler{},
	SlashReserveDevice: &modals.ReserveDeviceHandler{},
	SlashReleaseDevice: &modals.ReleaseDeviceHandler{},
}

type DeviceManager struct {
	Devices     *modals.DevicesInfo
	SlackClient *socketmode.Client
}

func (dm *DeviceManager) processEventApi(event socketmode.Event) {
	// The Event sent on the channel is not the same as the EventAPI events so we need to type cast it
	eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
	if !ok {
		log.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
		return
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
		log.Printf("Could not type cast the message to a Interaction callback: %v\n", interaction)
		return
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
		log.Printf("Could not type cast the message to a SlashCommand: %v\n", command)
		return
	}
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
	// Create a for loop that selects either the context cancellation or the events incomming
	for {
		select {
		// inscase context cancel is called exit the goroutine
		case <-ctx.Done():
			log.Println("Shutting down socketmode listener")
			return
		case event := <-dm.SlackClient.Events:
			// We have a new Events, let's type switch the event
			// Add more use cases here if you want to listen to other events.
			switch event.Type {
			// handle EventAPI events
			case socketmode.EventTypeEventsAPI:
				dm.processEventApi(event)
				// Handle Slash Commands
			case socketmode.EventTypeSlashCommand:
				dm.processSlashCommand(event)
				// Handle interaction events i.e. user voted in our poll etc.
			case socketmode.EventTypeInteractive:
				dm.processEventInteractive(event)
			}
		}
	}
}

// handleSlashCommand will take a slash command and route to the appropriate function
func (dm *DeviceManager) handleSlashCommand(command slack.SlashCommand) error {
	// TODO: Ignore commands from channels that the bot is not part of !!!
	handler, hasValue := SlashCommands[command.Command]
	if !hasValue {
		log.Printf("---> Unsupported command %s\n", command.Command)
		return nil
	}
	return dm.handleDeviceCommand(&command, handler)
}

func (dm *DeviceManager) handleDeviceCommand(
	command *slack.SlashCommand,
	handler ModalHandler,
) error {
	modalRequest := handler.GenerateModalRequest(command, *dm.Devices)
	_, err := dm.SlackClient.OpenView(command.TriggerID, modalRequest)
	if err != nil {
		return fmt.Errorf("Error opening view: %s", err)
	}
	return nil
}

func (dm *DeviceManager) handleInteractionEvent(interaction slack.InteractionCallback) error {
	switch interaction.Type {
	case slack.InteractionTypeViewSubmission:
		// NOTE: we use title text to determine which modal was submitted
		switch interaction.View.Title.Text {
		case modals.MReserveDeviceTitle:
			for _, selected := range interaction.View.State.Values[modals.MReserveDeviceActionId][modals.MReserveDeviceCheckboxId].SelectedOptions {
				for _, device := range *dm.Devices {
					if device.Name == selected.Value {
						device.Reserved = true
						device.ReservedBy = interaction.User.Name
						device.ReservedTime = time.Now()
					}
				}
			}
		case modals.MReleaseDeviceTitle:
			for _, selected := range interaction.View.State.Values[modals.MReleaseDeviceActionId][modals.MReleaseDeviceCheckboxId].SelectedOptions {
				for _, device := range *dm.Devices {
					if device.Name == selected.Value {
						device.Reserved = false
					}
				}
			}
		default:
		}
	default:

	}

	return nil
}
