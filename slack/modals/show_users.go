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

	return generateModalRequest("Testing status", allBlocks)
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

	var optionBlocks []*slack.OptionBlockObject

	sectionBlock := slack.NewOptionBlockObject(
		"name1",
		slack.NewTextBlockObject("plain_text", "Option1", false, false),
		slack.NewTextBlockObject("plain_text", "description1", false, false),
	)
	optionBlocks = append(optionBlocks, sectionBlock)

	sectionBlock2 := slack.NewOptionBlockObject(
		"name2",
		slack.NewTextBlockObject("plain_text", "Option2", false, false),
		slack.NewTextBlockObject("plain_text", "description2", false, false),
	)
	optionBlocks = append(optionBlocks, sectionBlock2)

	label1 := slack.NewTextBlockObject("plain_text", "Select", false, false)
	placeholder1 := slack.NewTextBlockObject("plain_text", "Select", false, false)
	element1 := slack.NewOptionsSelectBlockElement("static_select", placeholder1, "defd", optionBlocks...)

	selectionBlock := slack.NewInputBlock("input_block3", label1, element1)

	allBlocks = append(allBlocks, selectionBlock)

	usersMap, ok := usersM.(users.UserMap)
	if !ok {
		log.Fatalf("Expected DevicesInfo but got %v", usersM)
	}
	userSectionBlocks := generateSectionBlocks(usersMap)

	for _, user := range userSectionBlocks {
		// divSection := slack.NewDividerBlock()
		allBlocks = append(allBlocks, user)
		// allBlocks = append(allBlocks, divSection)
	}

	divSection1 := slack.NewDividerBlock()
	allBlocks = append(allBlocks, divSection1)

	add_user_field := slack.NewInputBlock("input_block1", slack.NewTextBlockObject("plain_text", "Add new user", false, false), slack.NewPlainTextInputBlockElement(slack.NewTextBlockObject("plain_text", "name.surname", false, false), ""))
	allBlocks = append(allBlocks, add_user_field)

	divSection := slack.NewDividerBlock()
	allBlocks = append(allBlocks, divSection)

	var userBlocks []*slack.OptionBlockObject

	for user_name, rights := range usersMap {
		var description string
		if rights == 1 {
			description = "Admin"
		} else {
			description = "Regular"
		}
		sectionBlock := slack.NewOptionBlockObject(
			user_name,
			slack.NewTextBlockObject("plain_text", user_name, false, false),
			slack.NewTextBlockObject("plain_text", description, false, false),
		)
		userBlocks = append(userBlocks, sectionBlock)
	}

	label := slack.NewTextBlockObject("plain_text", "Remove users", false, false)
	placeholder := slack.NewTextBlockObject("plain_text", "Select user", false, false)
	element := slack.NewOptionsSelectBlockElement("multi_static_select", placeholder, "asdf", userBlocks...)
	remove_user_field := slack.NewInputBlock("input_block2", label, element)

	allBlocks = append(allBlocks, remove_user_field)

	return allBlocks
}
