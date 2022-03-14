package handlers

import (
	"github.com/AngelVI13/slack-assistant/modals"
	"github.com/slack-go/slack"
)

type ModalHandler interface {
	GenerateModalRequest(modals.DevicesInfo) slack.ModalViewRequest
	GenerateBlocks(modals.DevicesInfo) []slack.Block
}
