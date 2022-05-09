package modals

import (
	"github.com/slack-go/slack"
)

const MAddUserTitle = "Add users"
const MAddUserActionId = "optionAction"
const MAddUserOptionId = "optionSelected"

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

	return allBlocks
}
