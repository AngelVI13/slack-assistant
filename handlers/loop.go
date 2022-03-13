package handlers

import (
	"context"
	"github.com/AngelVI13/slack-assistant/modals"
	"log"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type DeviceManager struct {
	DevicesInfo modals.DevicesInfo
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

	err := HandleInteractionEvent(interaction, dm.SlackClient)
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
	payload, err := HandleSlashCommand(command, dm.SlackClient)
	if err != nil {
		log.Fatal(err)
	}
	// Dont forget to acknowledge the request and send the payload
	// The payload is the response
	dm.SlackClient.Ack(*event.Request, payload)

}

// TODO: create a class SocketHandler that has all processing functions attached to it
//       and keeps a poitner to the socketClient & device Info data etc.
func (dm *DeviceManager) ProcessMessageLoop(ctx context.Context) {
	log.Println(dm.DevicesInfo)

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
