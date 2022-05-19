package slash

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/AngelVI13/slack-assistant/data"
	"github.com/AngelVI13/slack-assistant/users"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

const ResetReviewersCmd = "reset"
const ListReviewersCmd = "list"

// TODO: Add New methods for these structs that can take a pointer to DeviceManager and store
// all the stuff that might be needed (slack client, users data etc.) This will have to be done while
// processing each command, to make sure that DeviceManager is initialized by then
type ReviewHandler struct{}

func (h *ReviewHandler) Execute(command *slack.SlashCommand, slackClient *socketmode.Client, dataObj any) error {
	dataHolder, ok := dataObj.(*data.DataHolder)
	if !ok {
		log.Fatalf("Expected users data, but got something else: %+v", dataObj)
	}
	reviewersInfo := &dataHolder.Reviewers

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

	commandTxt := strings.TrimSpace(command.Text)

	// If instead of TaskID we have a subcommand -> handle that (only for admins)
	userInfo, ok := dataHolder.Users.Map[command.UserName]
	if ok && userInfo.Rights == users.ADMIN {
		if commandTxt == ResetReviewersCmd {
			resetReviewers(reviewersInfo, command, slackClient)
			return nil
		} else if commandTxt == ListReviewersCmd {
			listReviewers(reviewersInfo, command, slackClient)
			return nil
		}
	}

	// Otherwise, treat command text as TaskID and try to process that
	url, errorMsg := getTaskLink(commandTxt)
	if errorMsg != "" {
		slackClient.PostEphemeral(command.ChannelID, command.UserID, slack.MsgOptionText(errorMsg, false))
		return nil
	}

	reviewer := reviewersInfo.ChooseReviewer(command.UserName)
	reviewMsg := fmt.Sprintf("Reviewer for %s is <@%s>\n\n_Submitted by_: <@%s>\n_URL_: %s", commandTxt, reviewer.Id, command.UserID, url)

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

func resetReviewers(reviewersInfo *users.Reviewers, command *slack.SlashCommand, slackClient *socketmode.Client) {
	reviewersInfo.ResetCurrentReviewers()
	msg := fmt.Sprintf(
		"Reviewer list reset successfully! There are %d available reviewers.\n\n_For a full list of reviewers use *'/review %s'* command._",
		len(reviewersInfo.Current),
		ListReviewersCmd,
	)
	slackClient.PostEphemeral(command.ChannelID, command.UserID, slack.MsgOptionText(msg, false))
}

func listReviewers(reviewersInfo *users.Reviewers, command *slack.SlashCommand, slackClient *socketmode.Client) {
	reviewerNames := []string{}
	for _, r := range reviewersInfo.Current {
		reviewerNames = append(reviewerNames, r.Name)
	}
	msg := fmt.Sprintf(
		"Available reviewers (%d):\n\n\t* %s",
		len(reviewersInfo.Current),
		strings.Join(reviewerNames, "\n\t* "),
	)
	slackClient.PostEphemeral(command.ChannelID, command.UserID, slack.MsgOptionText(msg, false))
}
