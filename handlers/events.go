package handlers

import (
	"errors"
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// HandleEventMessage will take an event and handle it properly based on the type of event
func HandleEventMessage(event slackevents.EventsAPIEvent, client *socketmode.Client) error {
	switch event.Type {
	// First we check if this is an CallbackEvent
	case slackevents.CallbackEvent:

		innerEvent := event.InnerEvent
		// Yet Another Type switch on the actual Data to see if its an AppMentionEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			// The application has been mentioned since this Event is a Mention event
			// TODO: rework so that payload is returned and message posting happens on top level(here)
			err := HandleAppMentionEvent(ev, client)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupported event type")
	}
	return nil
}

// HandleAppMentionEvent is used to take care of the AppMentionEvent when the bot is mentioned
func HandleAppMentionEvent(event *slackevents.AppMentionEvent, client *socketmode.Client) error {
	// Grab the user name based on the ID of the one who mentioned the bot
	user, err := client.GetUserInfo(event.User)
	if err != nil {
		return err
	}

	// Create a help text from the supported slash commands
	// TODO: maybe add a command description?
	helpText := ""
	for key := range SlashCommands {
		helpText += fmt.Sprintf(":black_medium_small_square: %s\n", key)
	}
	// Send a message to the user
	text := fmt.Sprintf("Here are some useful commands, *%s*:\n%s", user.Name, helpText)

	// Send the message to the channel
	// The Channel is available in the event message
	_, _, err = client.PostMessage(event.Channel, slack.MsgOptionText(text, false))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
}
