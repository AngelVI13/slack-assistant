package modals

import (
	"log"

	"github.com/AngelVI13/slack-assistant/device"
	"github.com/slack-go/slack"
)

const MReleaseDeviceTitle = "Release Device"
const MReleaseDeviceActionId = "deviceSelected"
const MReleaseDeviceCheckboxId = "deviceCheckbox"

type ReleaseDeviceHandler struct{}

func (h *ReleaseDeviceHandler) GenerateModalRequest(data any) slack.ModalViewRequest {
	allBlocks := h.GenerateBlocks(data)
	return generateModalRequest(MReleaseDeviceTitle, allBlocks)
}

func (h *ReleaseDeviceHandler) GenerateBlocks(data any) []slack.Block {
	devices, ok := data.(device.DevicesInfo)
	if !ok {
		log.Fatal("Expected DevicesInfo but got something else")
	}
	deviceOptionBlocks := generateDeviceTakenOptionBlocks(devices)

	// If no devices are taken -> return a simple message to the user
	if len(deviceOptionBlocks) <= 0 {
		return []slack.Block{
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", "All devices are free, nothing to release", false, false),
				nil,
				nil,
			),
		}
	}
	// TODO: sort devices starting from devices I have reserved first
	//       so i don't have to search in the list which ones are mine

	// Turn device blocks to a poll/action element block
	deviceCheckboxGroup := slack.NewCheckboxGroupsBlockElement(MReleaseDeviceCheckboxId, deviceOptionBlocks...)
	actionBlocks := slack.NewActionBlock(MReleaseDeviceActionId, deviceCheckboxGroup)

	header := "Choose a device you would like to release"
	headerText := slack.NewTextBlockObject("mrkdwn", header, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Add header text and action(poll) elem to slice of modal blocks
	allBlocks := []slack.Block{headerSection, actionBlocks}
	return allBlocks
}
