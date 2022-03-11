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
	return generateInfoModalRequest(MShowDeviceTitle, allBlocks)
}

func (h *ShowDeviceHandler) GenerateBlocks(devices DevicesInfo) []slack.Block {
	deviceSectionBlocks := generateDeviceInfoBlocks(devices)

	var allBlocks []slack.Block
	for idx, device := range deviceSectionBlocks {
		divSection := slack.NewDividerBlock()
		allBlocks = append(allBlocks, device)

		// do not add separator after last element
		if idx < len(deviceSectionBlocks)-1 {
			allBlocks = append(allBlocks, divSection)
		}
	}
	return allBlocks
}
