package handlers

import (
	"context"
	"github.com/AngelVI13/slack-assistant/modals"
	"log"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func processEventApi(event socketmode.Event, client *slack.Client, socketClient *socketmode.Client) {
	// The Event sent on the channel is not the same as the EventAPI events so we need to type cast it
	eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
	if !ok {
		log.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
		return
	}
	// We need to send an Acknowledge to the slack server
	socketClient.Ack(*event.Request)
	// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch
	err := HandleEventMessage(eventsAPIEvent, client)
	if err != nil {
		// Replace with actual err handeling
		log.Fatal(err)
	}
}

func processEventInteractive(event socketmode.Event, client *slack.Client, socketClient *socketmode.Client) {
	interaction, ok := event.Data.(slack.InteractionCallback)
	if !ok {
		log.Printf("Could not type cast the message to a Interaction callback: %v\n", interaction)
		return
	}

	err := HandleInteractionEvent(interaction, client)
	if err != nil {
		log.Fatal(err)
	}
	socketClient.Ack(*event.Request)
}

func processSlashCommand(event socketmode.Event, client *slack.Client, socketClient *socketmode.Client) {

	// Just like before, type cast to the correct event type, this time a SlashEvent
	command, ok := event.Data.(slack.SlashCommand)
	if !ok {
		log.Printf("Could not type cast the message to a SlashCommand: %v\n", command)
		return
	}
	// handleSlashCommand will take care of the command
	payload, err := HandleSlashCommand(command, client)
	if err != nil {
		log.Fatal(err)
	}
	// Dont forget to acknowledge the request and send the payload
	// The payload is the response
	socketClient.Ack(*event.Request, payload)

}

// TODO: create a class SocketHandler that has all processing functions attached to it
//       and keeps a poitner to the socketClient & device Info data etc.
func ProcessMessageLoop(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {
	devicesInfo := modals.DevicesInfo{
		modals.DeviceInfo{"splinter", false},
		modals.DeviceInfo{"shredder", false},
		modals.DeviceInfo{"donatello", true},
	}
	log.Println(devicesInfo)

	// Create a for loop that selects either the context cancellation or the events incomming
	for {
		select {
		// inscase context cancel is called exit the goroutine
		case <-ctx.Done():
			log.Println("Shutting down socketmode listener")
			return
		case event := <-socketClient.Events:
			// We have a new Events, let's type switch the event
			// Add more use cases here if you want to listen to other events.
			switch event.Type {
			// handle EventAPI events
			case socketmode.EventTypeEventsAPI:
				processEventApi(event, client, socketClient)
				// Handle Slash Commands
			case socketmode.EventTypeSlashCommand:
				processSlashCommand(event, client, socketClient)
				// Handle interaction events i.e. user voted in our poll etc.
			case socketmode.EventTypeInteractive:
				processEventInteractive(event, client, socketClient)
			}
		}
	}
}
