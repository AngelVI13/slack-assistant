package modals

import (
	"fmt"
	"github.com/slack-go/slack"
)

const MUnauthorizedTitle = "Unauthorized user"

type UnauthorizedHandler struct{}

func (h *UnauthorizedHandler) GenerateModalRequest(command *slack.SlashCommand, devices DevicesInfo) slack.ModalViewRequest {
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

func (h *UnauthorizedHandler) GenerateBlocks(devices DevicesInfo) []slack.Block {
	return []slack.Block{}
}
