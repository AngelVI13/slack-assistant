package handlers

import (
	"fmt"
	"github.com/AngelVI13/slack-assistant/modals"
	"github.com/slack-go/slack"
	"time"
)

// HandleSlashCommand will take a slash command and route to the appropriate function
func HandleSlashCommand(command slack.SlashCommand, client *slack.Client) (interface{}, error) {
	// We need to switch depending on the command
	switch command.Command {
	case "/hello":
		// This was a hello command, so pass it along to the proper function
		return nil, HandleHelloCommand(command, client)
	case "/reserve-device":
		modalRequest := modals.GenerateReserveDeviceModalRequest(modals.Devices)
		_, err := client.OpenView(command.TriggerID, modalRequest)
		if err != nil {
			fmt.Printf("Error opening view: %s", err)
		}
		return nil, nil
	case "/release-device":
		modalRequest := modals.GenerateReleaseDeviceModalRequest(modals.Devices)
		_, err := client.OpenView(command.TriggerID, modalRequest)
		if err != nil {
			fmt.Printf("Error opening view: %s", err)
		}
		return nil, nil
	}

	return nil, nil
}

// HandleHelloCommand will take care of /hello submissions
func HandleHelloCommand(command slack.SlashCommand, client *slack.Client) error {
	// The Input is found in the text field so
	// Create the attachment and assigned based on the message
	attachment := slack.Attachment{}
	// Add Some default context like user who mentioned the bot
	attachment.Fields = []slack.AttachmentField{
		{
			Title: "Date",
			Value: time.Now().String(),
		}, {
			Title: "Initializer",
			Value: command.UserName,
		},
	}

	// Greet the user
	attachment.Text = fmt.Sprintf("Hello %s", command.Text)
	attachment.Color = "#4af030"

	// Send the message to the channel
	// The Channel is available in the command.ChannelID
	_, _, err := client.PostMessage(command.ChannelID, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
}

// HandleReserveDevice will trigger a Yes or No question to the initializer
func HandleReserveDevice(command slack.SlashCommand, client *slack.Client) (interface{}, error) {
	// Create the attachment and assigned based on the message
	attachment := slack.Attachment{}

	// Create the checkbox element
	checkbox := slack.NewCheckboxGroupsBlockElement("answer",
		slack.NewOptionBlockObject("splinter", &slack.TextBlockObject{Text: "Splinter", Type: slack.MarkdownType}, &slack.TextBlockObject{Text: "Port: 5568", Type: slack.MarkdownType}),
		slack.NewOptionBlockObject("shredder", &slack.TextBlockObject{Text: "Shredder", Type: slack.MarkdownType}, &slack.TextBlockObject{Text: "Port: 5555", Type: slack.MarkdownType}),
	)
	// Create the Accessory that will be included in the Block and add the checkbox to it
	accessory := slack.NewAccessory(checkbox)
	// Add Blocks to the attachment
	attachment.Blocks = slack.Blocks{
		BlockSet: []slack.Block{
			// Create a new section block element and add some text and the accessory to it
			slack.NewSectionBlock(
				&slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: "Which device would you like to reserve?",
				},
				nil,
				accessory,
			),
		},
	}

	// TODO: what do the following properties do? Can't see this in slack
	attachment.Text = "Rate the tutorial"
	attachment.Color = "#4af030"
	return attachment, nil
}
