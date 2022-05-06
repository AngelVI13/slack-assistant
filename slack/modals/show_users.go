package modals

import (
	"fmt"
	"log"

	"github.com/AngelVI13/slack-assistant/utils/users"
	"github.com/slack-go/slack"
)

const MShowUsersTitle = "Show users"
const MShowUsersActionId = "optionAction"
const MShowUsersOptionId = "optionSelected"

type ShowUsersHandler struct{}

func (h *ShowUsersHandler) GenerateModalRequest(users any) slack.ModalViewRequest {
	allBlocks := h.GenerateBlocks(users)

	return generateModalRequest(MShowUsersTitle, allBlocks)
}

func generateSectionBlocks(users users.UserMap) []*slack.SectionBlock {
	var userBlocks []*slack.SectionBlock

	for user, rights := range users {
		text := fmt.Sprintf("%s, %v", user, rights)
		sectionText := slack.NewTextBlockObject("plain_text", text, false, false)
		sectionBlock := slack.NewSectionBlock(sectionText, nil, nil)
		userBlocks = append(userBlocks, sectionBlock)
	}

	return userBlocks
}

func (h *ShowUsersHandler) GenerateBlocks(usersM any) []slack.Block {
	var allBlocks []slack.Block

	usersMap, ok := usersM.(users.UserMap)
	if !ok {
		log.Fatalf("Expected DevicesInfo but got %v", usersM)
	}
	userSectionBlocks := generateSectionBlocks(usersMap)

	for _, user := range userSectionBlocks {
		allBlocks = append(allBlocks, user)
	}

	return allBlocks
}
