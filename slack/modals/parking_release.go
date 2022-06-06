package modals

import (
	"fmt"
	"log"

	"github.com/AngelVI13/slack-assistant/parking"
	"github.com/slack-go/slack"
)

const (
	MParkingReleaseTitle = "Temporary release a parking spot"

	// TODO: Add 2 buttons for Release for Special users (on the booking page)
	//       1. Button for temporary release of spot -> leads to this modal
	//       2. Button for permament release (acts the same as release for non-special users)
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
	allBlocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("Temporarily release space: %d (%d floor)", space.Number, space.Floor), false, false),
			nil,
			nil,
		),
		slack.NewActionBlock(
			ReleaseBlockId,
			slack.NewDatePickerBlockElement(ReleaseStartDateActionId),
			slack.NewDatePickerBlockElement(ReleaseEndDateActionId),
		),
	}
	return allBlocks
}
