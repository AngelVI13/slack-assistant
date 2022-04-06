package handlers

import (
	"github.com/slack-go/slack"
)

type ModalHandler interface {
	GenerateModalRequest(*slack.SlashCommand, any) slack.ModalViewRequest
	GenerateBlocks(any) []slack.Block
}
