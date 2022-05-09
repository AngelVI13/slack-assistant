package modals

import (
	"fmt"

	"github.com/AngelVI13/slack-assistant/utils/users"
	"github.com/slack-go/slack"
)

const MRemoveUsersTitle = "Remove users"
const MRemoveUsersActionId = "MRemoveUsersActionId"
const MRemoveUsersOptionId = "MRemoveUsersOptionId"

type RemoveUsersHandler struct{}

func (h *RemoveUsersHandler) GenerateModalRequest(users any) slack.ModalViewRequest {
	allBlocks := h.GenerateBlocks(users)

	return generateModalRequest(MRemoveUsersTitle, allBlocks)
}

func (h *RemoveUsersHandler) GenerateBlocks(usersM any) []slack.Block {
	var allBlocks []slack.Block

	var userBlocks []*slack.OptionBlockObject

	usersMap := usersM.(users.UserMap)
	for user, rights := range usersMap {
		sectionBlock := slack.NewOptionBlockObject(
			user,
			slack.NewTextBlockObject("plain_text", user, false, false),
			slack.NewTextBlockObject("plain_text", fmt.Sprintf("%v", rights), false, false),
		)
		userBlocks = append(userBlocks, sectionBlock)
	}

	label := slack.NewTextBlockObject("plain_text", "Remove users", false, false)
	placeholder := slack.NewTextBlockObject("plain_text", "Select user", false, false)
	element := slack.NewOptionsSelectBlockElement("multi_static_select", placeholder, MRemoveUsersOptionId, userBlocks...)
	remove_user_field := slack.NewInputBlock(MRemoveUsersActionId, label, element)
	allBlocks = append(allBlocks, remove_user_field)

	return allBlocks
}
