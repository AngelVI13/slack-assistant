package modals

import (
	"github.com/slack-go/slack"
)

const MReleaseDeviceTitle = "Release Device"
const MReleaseDeviceActionId = "deviceSelected"
const MReleaseDeviceCheckboxId = "deviceCheckbox"

type ReleaseDeviceHandler struct {}

func (h *ReleaseDeviceHandler) GenerateModalRequest(devices DevicesInfo) slack.ModalViewRequest {
    allBlocks := h.GenerateBlocks(devices)
	return generateModalRequest(MReleaseDeviceTitle, allBlocks)
}
    
func (h *ReleaseDeviceHandler) GenerateBlocks(devices DevicesInfo) []slack.Block {
	deviceOptionBlocks := generateDeviceOptionBlocks(devices, getTakenDevices)
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
