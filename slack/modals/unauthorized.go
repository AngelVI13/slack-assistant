package modals

import (
	"fmt"
	"github.com/slack-go/slack"
	"github.com/AngelVI13/slack-assistant/device"
)

const MUnauthorizedTitle = "Unauthorized user"

type UnauthorizedHandler struct{}

func (h *UnauthorizedHandler) GenerateModalRequest(command *slack.SlashCommand, devices device.DevicesInfo) slack.ModalViewRequest {
	userMsg := fmt.Sprintf(
		":bust_in_silhouette: *%s*,\n\nYou don't have access to execute command *%s*. This incident will be reported!",
		command.UserName,
		command.Command)
	allBlocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", userMsg, false, false),
			nil,
			nil,
		),
	}
	return generateModalRequest(MUnauthorizedTitle, allBlocks)
}

func (h *UnauthorizedHandler) GenerateBlocks(devices device.DevicesInfo) []slack.Block {
	return []slack.Block{}
}
