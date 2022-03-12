package handlers

import (
	"errors"
	"fmt"
	"github.com/AngelVI13/slack-assistant/modals"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"log"
	"strings"
	"time"
)

// HandleEventMessage will take an event and handle it properly based on the type of event
func HandleEventMessage(event slackevents.EventsAPIEvent, client *slack.Client) error {
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
func HandleAppMentionEvent(event *slackevents.AppMentionEvent, client *slack.Client) error {

	// Grab the user name based on the ID of the one who mentioned the bot
	user, err := client.GetUserInfo(event.User)
	if err != nil {
		return err
	}
	// Check if the user said Hello to the bot
	text := strings.ToLower(event.Text)

	// Create the attachment and assigned based on the message
	attachment := slack.Attachment{}
	// Add Some default context like user who mentioned the bot
	attachment.Fields = []slack.AttachmentField{
		{
			Title: "Date",
			Value: time.Now().String(),
		}, {
			Title: "Initializer",
			Value: user.Name,
		},
	}
	if strings.Contains(text, "hello") {
		// Greet the user
		attachment.Text = fmt.Sprintf("Hello %s", user.Name)
		attachment.Pretext = "Greetings"
		attachment.Color = "#4af030"
	} else {
		// Send a message to the user
		attachment.Text = fmt.Sprintf("How can I help you %s?", user.Name)
		attachment.Pretext = "How can I be of service"
		attachment.Color = "#3d3d3d"
	}
	// Send the message to the channel
	// The Channel is available in the event message
	_, _, err = client.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
}

func HandleInteractionEvent(interaction slack.InteractionCallback, client *slack.Client) error {
	// This is where we would handle the interaction
	// Switch depending on the Type
	log.Printf("The action called is: %s\n", interaction.ActionID)
	log.Printf("The response was of type: %s\n", interaction.Type)
	switch interaction.Type {
	case slack.InteractionTypeBlockActions:
		// This is a block action, so we need to handle it

		for _, action := range interaction.ActionCallback.BlockActions {
			log.Printf("%+v", action)
			log.Println("Selected option: ", action.SelectedOptions)

		}

	case slack.InteractionTypeViewSubmission:
		// NOTE: we can use title text to determine which modal was submitted
		log.Printf("----> %+v", interaction.View.Title.Text == modals.MReserveDeviceTitle)
		for _, selected := range interaction.View.State.Values[modals.MReserveDeviceActionId][modals.MReserveDeviceCheckboxId].SelectedOptions {
			log.Printf("%+v\n", selected.Value)
		}
	default:

	}

	return nil
}
