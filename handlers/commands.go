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
		return nil, HandleHelloCommand(command, client)
	case "/reserve-device":
		return nil, HandleDeviceCommand(command, client, modals.Devices, modals.GenerateReserveDeviceModalRequest)
	case "/release-device":
		return nil, HandleDeviceCommand(command, client, modals.Devices, modals.GenerateReleaseDeviceModalRequest)
	case "/show-devices":
		return nil, HandleDeviceCommand(command, client, modals.Devices, modals.GenerateShowDeviceModalRequest)
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

func HandleDeviceCommand(
	command slack.SlashCommand,
	client *slack.Client,
	devicesInfo modals.DevicesInfo,
	generateModal func(modals.DevicesInfo) slack.ModalViewRequest,
) error {
	modalRequest := generateModal(devicesInfo)
	_, err := client.OpenView(command.TriggerID, modalRequest)
	if err != nil {
		return fmt.Errorf("Error opening view: %s", err)
	}
	return nil
}
