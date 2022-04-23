package slash

import "github.com/slack-go/slack"

type SlashHandler interface {
	Execute(command *slack.SlashCommand, data any)
}
