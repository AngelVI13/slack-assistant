package modals

import (
	"github.com/slack-go/slack"
)

const MReserveDeviceTitle = "Reserve Device"
const MReserveDeviceActionId = "deviceSelected"
const MReserveDeviceCheckboxId = "deviceCheckbox"

type ReserveDeviceHandler struct {}

func (h *ReserveDeviceHandler) GenerateModalRequest(devices DevicesInfo) slack.ModalViewRequest {
    allBlocks := h.GenerateBlocks(devices)
	return generateModalRequest(MReserveDeviceTitle, allBlocks)
}
    
func (h *ReserveDeviceHandler) GenerateBlocks(devices DevicesInfo) []slack.Block {
	deviceOptionBlocks := generateDeviceOptionBlocks(devices, getFreeDevices)
	// Turn device blocks to a poll/action element block
	deviceCheckboxGroup := slack.NewCheckboxGroupsBlockElement(MReserveDeviceCheckboxId, deviceOptionBlocks...)
	actionBlocks := slack.NewActionBlock(MReserveDeviceActionId, deviceCheckboxGroup)

	header := "Choose a device you would like to reserve"
	headerText := slack.NewTextBlockObject("mrkdwn", header, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Add header text and action(poll) elem to slice of modal blocks
	allBlocks := []slack.Block{headerSection, actionBlocks}
    return allBlocks
}

