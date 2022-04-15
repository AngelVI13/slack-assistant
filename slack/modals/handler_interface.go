package modals

import (
	"github.com/slack-go/slack"
)

type ModalHandler interface {
	GenerateModalRequest(any) slack.ModalViewRequest
	GenerateBlocks(any) []slack.Block
	ChangeAction(action string)
}
