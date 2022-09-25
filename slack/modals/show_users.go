package modals

import (
	"fmt"
	"log"

	"github.com/AngelVI13/slack-assistant/users"
	"github.com/slack-go/slack"
)

const MShowUsersTitle = "Show users"
const MShowUsersActionId = "optionAction"
const MShowUsersOptionId = "optionSelected"

type ShowUsersHandler struct{}

func (h *ShowUsersHandler) GenerateModalRequest(command *slack.SlashCommand, users ...any) slack.ModalViewRequest {
	allBlocks := h.GenerateBlocks(command, users...)

	return generateModalRequest(MShowUsersTitle, allBlocks)
}

func generateSectionBlocks(users users.UsersMap) []*slack.SectionBlock {
	var userBlocks []*slack.SectionBlock

	for user, rights := range users {
		text := fmt.Sprintf("%s, %v", user, rights)
		sectionText := slack.NewTextBlockObject("plain_text", text, false, false)
		sectionBlock := slack.NewSectionBlock(sectionText, nil, nil)
		userBlocks = append(userBlocks, sectionBlock)
	}

	return userBlocks
}

func (h *ShowUsersHandler) GenerateBlocks(command *slack.SlashCommand, usersM ...any) []slack.Block {
	var allBlocks []slack.Block

	if len(usersM) != 1 {
		log.Fatal("Incorrect num of params for ShowUsersHandler: 1")
	}
	rawUsers := usersM[0]

	usersMap, ok := rawUsers.(users.UsersMap)
	if !ok {
		log.Fatalf("Expected DevicesInfo but got %v", usersM)
	}
	userSectionBlocks := generateSectionBlocks(usersMap)

	for _, user := range userSectionBlocks {
		allBlocks = append(allBlocks, user)
	}

	return allBlocks
}
