package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AngelVI13/slack-assistant/data"
	"github.com/AngelVI13/slack-assistant/slack/modals"
	"github.com/AngelVI13/slack-assistant/slack/slash"

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
	"/parking": modals.NewCustomOptionModalHandler(
		modals.ParkingActionMap,
		modals.DefaultParkingAction,
		modals.ParkingModalInfo,
	),
}

var SlashCommandsForHandlers = map[string]slash.SlashHandler{
	"/review": &slash.ReviewHandler{},
}

type OptionModalData struct {
	Handler modals.OptionModalHandler
	Command *slack.SlashCommand
}

type SlackBot struct {
	Data *data.DataHolder

	SlackClient *socketmode.Client
	// Whenever we are dealing with a modal that contains a state switching option
	// keep a pointer to it so we can change states
	CurrentOptionModalData *OptionModalData
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
	modalRequest := handler.GenerateModalRequest(command, bot.Data.Devices.GetDevicesInfo(command.UserName))

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
		return handler.Execute(&command, bot.SlackClient, bot.Data)
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
		data = bot.Data.ParkingLot.GetSpacesInfo(command.UserName)
	} else {
		data = bot.Data.Devices.GetDevicesInfo(command.UserName)
	}

	// In case we are dealing with an OptionModalHandler save pointer to it
	// so we can change its state when needed
	optionHandler, ok := handler.(modals.OptionModalHandler)
	if ok {
		bot.CurrentOptionModalData = &OptionModalData{
			Handler: optionHandler,
			Command: command,
		}
		bot.CurrentOptionModalData.Handler.Reset()
	} else {
		bot.CurrentOptionModalData = nil
	}

	modalRequest := handler.GenerateModalRequest(command, data)
	_, err := bot.SlackClient.OpenView(command.TriggerID, modalRequest)
	if err != nil {
		return fmt.Errorf("Error opening view: %s", err)
	}
	return nil
}

func (bot *SlackBot) handleInteractionEvent(interaction slack.InteractionCallback) error {
	switch interaction.Type {
	case slack.InteractionTypeViewSubmission:
		bot.handleViewSubmission(&interaction)
	case slack.InteractionTypeBlockActions:
	default:

	}

	return nil
}

func (bot *SlackBot) handleViewSubmission(interaction *slack.InteractionCallback) {
	// NOTE: we use title text to determine which modal was submitted
	switch interaction.View.Title.Text {
	case modals.MRestartProxyTitle:
		restartProxySubmission(bot.SlackClient, interaction, bot.Data.Devices)
	case modals.MRemoveUsersTitle:
		removeUserSubmission(interaction, bot.Data.Users)
	case modals.MAddUserTitle:
		addUserSubmission(bot.SlackClient, interaction, bot.Data.Users, bot.Data.Reviewers)
	default:
	}
}

func (bot *SlackBot) handleBlockActions(interaction *slack.InteractionCallback) {
	if bot.CurrentOptionModalData.Handler == nil {
		log.Fatalf(
			`Did not have a valid pointer to OptionModal,
                    please make sure to close any open modals before restarting the bot`,
		)
	}

	var updatedView *slack.ModalViewRequest

	switch interaction.View.Title.Text {
	case modals.MDeviceTitle:
		updatedView = handleDeviceActions(bot, interaction)
	case modals.MShowUsersTitle, modals.MRemoveUsersTitle, modals.MAddUserTitle:
		updatedView = handleUserActions(bot, interaction)
	case modals.MParkingTitle:
		updatedView = handleParkingActions(bot, interaction)
	case modals.MParkingBookingTitle: // TODO: Why is this not the parking release title instead?
		updatedView = handleParkingBooking(bot, interaction)
	default:
	}

	// Update view if a handler generated an update
	if updatedView != nil {
		_, err := bot.SlackClient.UpdateView(*updatedView, "", "", interaction.View.ID)
		if err != nil {
			log.Fatal(err)
		}
	}
}
