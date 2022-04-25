package slash

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

// TODO: Add New methods for these structs that can take a pointer to DeviceManager and store
// all the stuff that might be needed (slack client, users data etc.) This will have to be done while
// processing each command, to make sure that DeviceManager is initialized by then
type ReviewHandler struct{}

func (h *ReviewHandler) Execute(command *slack.SlashCommand, slackClient *socketmode.Client, data any) error {
	taskId := strings.TrimSpace(command.Text)
	var url string
	if len(strings.Split(taskId, " ")) > 1 {
		// TODO: maybe open a modal with this message ?
		msg := fmt.Sprintf("Incorrect task ID: *'%s'*.\nTask ID should be of the format 4AP2-1234 (for polarion) or 1234 (for azure PR)", taskId)
		slackClient.PostEphemeral(command.UserID, command.UserID, slack.MsgOptionText(msg, false))
		return nil
	} else if strings.Contains(taskId, "4AP2-") {
		url = fmt.Sprintf("https://alm-machines001.schweinfurt.germany.fresenius.de/polarion/#/project/4008APackage2/workitem?id=%s", taskId)
	} else if _, err := strconv.Atoi(taskId); err == nil {
		url = fmt.Sprintf("https://dev.azure.com/FMC-SSM/TestAutomation/_git/TestAutomation/pullrequest/%s", taskId)
	}

	userIdPlaceholder := "U9K74SZT7" // currently this is AI ID
	reviewMsg := fmt.Sprintf("Reviewer for %s is <@%s>\n\n _URL_ \n%s", taskId, userIdPlaceholder, url)
	slackClient.PostEphemeral(command.UserID, command.UserID, slack.MsgOptionText(reviewMsg, false))
	return nil
}
