package modals

import (
	"fmt"
	"log"

	"github.com/AngelVI13/slack-assistant/parking"
	"github.com/slack-go/slack"
)

const (
	MParkingBookingTitle = "Parking Booking"

	ReserveParkingActionId         = "reserve"
	ReserveParkingWithAutoActionId = "withAuto"
	ReleaseParkingActionId         = "release"
)

type ParkingBookingHandler struct{}

func (h *ParkingBookingHandler) GenerateModalRequest(command *slack.SlashCommand, data any) slack.ModalViewRequest {
	allBlocks := h.GenerateBlocks(command, data)
	return generateInfoModalRequest(MParkingBookingTitle, allBlocks)
}

func (h *ParkingBookingHandler) GenerateBlocks(command *slack.SlashCommand, data any) []slack.Block {
	spaces, ok := data.(parking.SpacesInfo)
	if !ok {
		log.Fatal("Expected SpacesInfo but got something else")
	}
	spacesSectionBlocks := generateParkingInfoBlocks(spaces)

	allBlocks := spacesSectionBlocks
	return allBlocks
}

// generateParkingInfo Generate sections of text that contain device info such as status (taken/free), ip, port, taken by etc..
func generateParkingInfo(spaces parking.SpacesInfo) []slack.SectionBlock {
	var sections []slack.SectionBlock
	for _, space := range spaces {
		status := space.GetStatusDescription()
		emoji := space.GetStatusEmoji()

		spaceProps := space.GetPropsText()
		text := fmt.Sprintf("%s *%s* \t%s\n\t\t%s", emoji, fmt.Sprint(space.Number), spaceProps, status)
		sectionText := slack.NewTextBlockObject("mrkdwn", text, false, false)
		sectionBlock := slack.NewSectionBlock(sectionText, nil, nil)

		sections = append(sections, *sectionBlock)
	}
	return sections
}

func generateParkingButtons(space *parking.ParkingSpace) []slack.BlockElement {
	var buttons []slack.BlockElement

	if space.Reserved {
		releaseButton := slack.NewButtonBlockElement(
			ReleaseParkingActionId,
			fmt.Sprint(space.Number),
			slack.NewTextBlockObject("plain_text", "Release!", true, false),
		)
		releaseButton = releaseButton.WithStyle(slack.StyleDanger)
		buttons = append(buttons, releaseButton)
	} else {

		actionButtonText := "Reserve!"
		reserveWithAutoButton := slack.NewButtonBlockElement(
			ReserveWithAutoActionId,
			fmt.Sprint(space.Number),
			slack.NewTextBlockObject("plain_text", fmt.Sprintf("%s :eject:", actionButtonText), true, false),
		)
		reserveWithAutoButton = reserveWithAutoButton.WithStyle(slack.StylePrimary)
		buttons = append(buttons, reserveWithAutoButton)

		reserveButton := slack.NewButtonBlockElement(ReserveParkingActionId, fmt.Sprint(space.Number), slack.NewTextBlockObject("plain_text", actionButtonText, true, false))
		buttons = append(buttons, reserveButton)
	}
	return buttons
}

// generateParkingInfoBlocks Generates device block objects to be used as elements in modal
func generateParkingInfoBlocks(spaces parking.SpacesInfo) []slack.Block {
	div := slack.NewDividerBlock()
	parkingSpaceSections := generateParkingInfo(spaces)

	var parkingSpaceBlocks []slack.Block
	for idx, device := range spaces {
		sectionBlock := parkingSpaceSections[idx]
		buttons := generateParkingButtons(device)

		actions := slack.NewActionBlock("", buttons...)
		parkingSpaceBlocks = append(parkingSpaceBlocks, sectionBlock, actions, div)
	}

	return parkingSpaceBlocks
}
