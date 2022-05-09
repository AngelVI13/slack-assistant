package slash

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/AngelVI13/slack-assistant/users"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

// TODO: Add New methods for these structs that can take a pointer to DeviceManager and store
// all the stuff that might be needed (slack client, users data etc.) This will have to be done while
// processing each command, to make sure that DeviceManager is initialized by then
type ReviewHandler struct{}

func (h *ReviewHandler) Execute(command *slack.SlashCommand, slackClient *socketmode.Client, data any) error {
	taskId := strings.TrimSpace(command.Text)
	url, errorMsg := getTaskLink(taskId)
	if errorMsg != "" {
		slackClient.PostEphemeral(command.UserID, command.UserID, slack.MsgOptionText(errorMsg, false))
		return nil
	}

	reviewersInfo, ok := data.(*users.Reviewers)
	if !ok {
		log.Fatalf("Expected users data, but got something else: %+v", data)
	}

	// If command is invoked from somewhere else than the required channel -> raise error
	if reviewersInfo.ChannelId != command.ChannelID {
		usedCommand := fmt.Sprintf("%s %s", command.Command, command.Text)
		errorMsg := fmt.Sprintf("Review command `%s` must be used inside channel <#%s>!", usedCommand, reviewersInfo.ChannelId)

		slackClient.PostEphemeral(
			command.UserID,
			command.UserID,
			slack.MsgOptionText(errorMsg, false),
		)
		return nil
	}

	reviewer := reviewersInfo.ChooseReviewer(command.UserName)
	reviewMsg := fmt.Sprintf("Reviewer for %s is <@%s>\n\n_Submitted by_: <@%s>\n_URL_: %s", taskId, reviewer.Id, command.UserID, url)

	// NOTE: the bot must be present in the channel otherwise, no response will be visible
	slackClient.PostMessage(reviewersInfo.ChannelId, slack.MsgOptionText(reviewMsg, false))
	return nil
}

func getTaskLink(taskId string) (url string, errorMsg string) {
	incorrectTaskId := fmt.Sprintf("Incorrect task ID: *'%s'*.\nTask ID should be of the format 4AP2-1234 (for polarion) or 1234 (for azure PR)", taskId)

	if len(strings.Split(taskId, " ")) > 1 {
		errorMsg = incorrectTaskId
	} else if strings.Contains(taskId, "4AP2-") {
		url = fmt.Sprintf("https://alm-machines001.schweinfurt.germany.fresenius.de/polarion/#/project/4008APackage2/workitem?id=%s", taskId)
	} else if _, err := strconv.Atoi(taskId); err == nil {
		url = fmt.Sprintf("https://dev.azure.com/FMC-SSM/TestAutomation/_git/TestAutomation/pullrequest/%s", taskId)
	} else {
		errorMsg = incorrectTaskId
	}

	return url, errorMsg
}
