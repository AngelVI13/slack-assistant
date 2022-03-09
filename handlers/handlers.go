package handlers

import (
	"github.com/AngelVI13/slack-assistant/modals"
	"github.com/slack-go/slack"
)

type ModalHandler interface {
	// HandleSlashCommand(slack.SlashCommand, *slack.Client, modals.DevicesInfo)
	GenerateModalRequest(modals.DevicesInfo) slack.ModalViewRequest
	GenerateBlocks(modals.DevicesInfo) []slack.Block
}
