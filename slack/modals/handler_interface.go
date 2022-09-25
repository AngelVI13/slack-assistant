package modals

import (
	"github.com/slack-go/slack"
)

type (
	ModalHandler interface {
		GenerateModalRequest(*slack.SlashCommand, ...any) slack.ModalViewRequest
		GenerateBlocks(*slack.SlashCommand, ...any) []slack.Block
	}

	OptionModalHandler interface {
		ModalHandler
		ChangeAction(action string) bool
		Reset()
	}
)
