package handlers

import (
	"github.com/AngelVI13/slack-assistant/device"
	"github.com/slack-go/slack"
)

type ModalHandler interface {
	GenerateModalRequest(*slack.SlashCommand, device.DevicesInfo) slack.ModalViewRequest
	GenerateBlocks(device.DevicesInfo) []slack.Block
}
