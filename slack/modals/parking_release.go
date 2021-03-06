package modals

import (
	"fmt"
	"log"

	"github.com/AngelVI13/slack-assistant/parking"
	"github.com/slack-go/slack"
)

const (
	MParkingReleaseTitle = "Temporary release a parking spot"

	TempReleaseParkingActionId = "tempRelease"
	ReleaseStartDateActionId   = "releaseStartDate"
	ReleaseEndDateActionId     = "releaseEndDate"
	ReleaseBlockId             = "releaseBlockId"
)

type ParkingReleaseHandler struct{}

func (h *ParkingReleaseHandler) GenerateModalRequest(command *slack.SlashCommand, data any) slack.ModalViewRequest {
	allBlocks := h.GenerateBlocks(command, data)
	return generateInfoModalRequest(MParkingBookingTitle, allBlocks)
}

func (h *ParkingReleaseHandler) GenerateBlocks(command *slack.SlashCommand, data any) []slack.Block {
	space, ok := data.(*parking.ParkingSpace)
	if !ok {
		log.Fatal("Expected ParkingSpace but got something else")
	}
	description := slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Temporarily release space: %d (%d floor)", space.Number, space.Floor), false, false),
		nil,
		nil,
	)

	startDate := slack.NewDatePickerBlockElement(ReleaseStartDateActionId)
	startDate.Placeholder = slack.NewTextBlockObject("plain_text", "Select START date", false, false)

	endDate := slack.NewDatePickerBlockElement(ReleaseEndDateActionId)
	endDate.Placeholder = slack.NewTextBlockObject("plain_text", "Select END date", false, false)

	calendarsSection := slack.NewActionBlock(
		ReleaseBlockId,
		startDate,
		endDate,
	)

	allBlocks := []slack.Block{
		description,
		calendarsSection,
	}
	return allBlocks
}
