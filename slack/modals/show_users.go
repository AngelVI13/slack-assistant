package modals

import (
	"fmt"
	"log"

	"github.com/AngelVI13/slack-assistant/utils/users"
	"github.com/slack-go/slack"
)

type ShowUsersHandler struct{}

func (h *ShowUsersHandler) GenerateModalRequest(command *slack.SlashCommand, users any) slack.ModalViewRequest {
	allBlocks := h.GenerateBlocks(users)
	return generateInfoModalRequest("Testing status", allBlocks)
}

func generateSectionBlocks(users users.UserMap) []*slack.SectionBlock {
	var userBlocks []*slack.SectionBlock

	for user, rights := range users {

		text := fmt.Sprintf("%s, %v", user, rights)
		sectionText := slack.NewTextBlockObject("mrkdwn", text, false, false)
		sectionBlock := slack.NewSectionBlock(sectionText, nil, nil)

		userBlocks = append(userBlocks, sectionBlock)
	}

	return userBlocks
}

func (h *ShowUsersHandler) GenerateBlocks(usersM any) []slack.Block {
	usersMap, ok := usersM.(users.UserMap)
	if !ok {
		log.Fatalf("Expected DevicesInfo but got %v", usersM)
	}
	userSectionBlocks := generateSectionBlocks(usersMap)

	var allBlocks []slack.Block
	for idx, user := range userSectionBlocks {
		divSection := slack.NewDividerBlock()
		allBlocks = append(allBlocks, user)

		// do not add separator after last element
		if idx < len(userSectionBlocks)-1 {
			allBlocks = append(allBlocks, divSection)
		}
	}
	return allBlocks
}
