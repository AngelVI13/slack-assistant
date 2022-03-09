package modals

import (
	"github.com/slack-go/slack"
)

const MShowDeviceTitle = "Device Status"
const MShowDeviceActionId = "deviceSelected"
const MShowDeviceCheckboxId = "deviceCheckbox"

type ShowDeviceHandler struct{}

func (h *ShowDeviceHandler) GenerateModalRequest(devices DevicesInfo) slack.ModalViewRequest {
	allBlocks := h.GenerateBlocks(devices)
	return generateModalRequest(MShowDeviceTitle, allBlocks)
}

func (h *ShowDeviceHandler) GenerateBlocks(devices DevicesInfo) []slack.Block {
	deviceOptionBlocks := generateDeviceOptionBlocks(devices, getAllDevices)
	// Turn device blocks to a poll/action element block
	deviceCheckboxGroup := slack.NewCheckboxGroupsBlockElement(MShowDeviceCheckboxId, deviceOptionBlocks...)
	actionBlocks := slack.NewActionBlock(MShowDeviceActionId, deviceCheckboxGroup)

	header := "Devices status"
	headerText := slack.NewTextBlockObject("mrkdwn", header, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Add header text and action(poll) elem to slice of modal blocks
	allBlocks := []slack.Block{headerSection, actionBlocks}
	return allBlocks
}
