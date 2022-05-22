package modals

import (
	"github.com/slack-go/slack"
)

const MAddUserTitle = "Add users"
const MAddUserActionId = "optionAction"
const MAddUserOptionId = "optionSelected"
const MAddUserAccessRightActionId = "optionActionAccessRight"
const MAddUserAccessRightOptionId = "optionSelectedAccessRight"
const MAddUserAccessRightOption = "Access right option"
const MAddUserReviewerOption = "Reviewer option"

type AddUserHandler struct{}

func (h *AddUserHandler) GenerateModalRequest(users any) slack.ModalViewRequest {
	allBlocks := h.GenerateBlocks(users)

	return generateModalRequest(MAddUserTitle, allBlocks)
}

func (h *AddUserHandler) GenerateBlocks(usersM any) []slack.Block {

	var allBlocks []slack.Block

	label := slack.NewTextBlockObject("plain_text", "Add users", false, false)
	placeholder := slack.NewTextBlockObject("plain_text", "Select users", false, false)
	element := slack.NewOptionsSelectBlockElement("multi_users_select", placeholder, MAddUserOptionId)
	add_user_field := slack.NewInputBlock(MAddUserActionId, label, element)
	allBlocks = append(allBlocks, add_user_field)

	var sectionBlocks []*slack.OptionBlockObject

	sectionBlock1 := slack.NewOptionBlockObject(
		MAddUserAccessRightOption,
		slack.NewTextBlockObject("mrkdwn", "Admin", false, false),
		slack.NewTextBlockObject("mrkdwn", "Select to assign Admin rights.", false, false),
	)
	sectionBlock2 := slack.NewOptionBlockObject(
		MAddUserReviewerOption,
		slack.NewTextBlockObject("mrkdwn", "Reviewer", false, false),
		slack.NewTextBlockObject("mrkdwn", "Select to assign Reviewer option.", false, false),
	)

	sectionBlocks = append(sectionBlocks, sectionBlock1, sectionBlock2)

	deviceCheckboxGroup := slack.NewCheckboxGroupsBlockElement(MAddUserAccessRightOptionId, sectionBlocks...)
	actionBlock := slack.NewActionBlock(MAddUserAccessRightActionId, deviceCheckboxGroup)
	allBlocks = append(allBlocks, actionBlock)

	return allBlocks
}
