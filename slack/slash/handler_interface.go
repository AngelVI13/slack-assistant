package slash

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type SlashHandler interface {
	Execute(command *slack.SlashCommand, slackClient *socketmode.Client, data any) error
}
