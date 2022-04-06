package modals

import (
	"log"

	"github.com/AngelVI13/slack-assistant/device"
	"github.com/slack-go/slack"
)

const MShowDeviceTitle = "Device Status"
const MShowDeviceActionId = "deviceSelected"
const MShowDeviceCheckboxId = "deviceCheckbox"

type ShowDeviceHandler struct{}

func (h *ShowDeviceHandler) GenerateModalRequest(command *slack.SlashCommand, data any) slack.ModalViewRequest {
	allBlocks := h.GenerateBlocks(data)
	return generateInfoModalRequest(MShowDeviceTitle, allBlocks)
}

func (h *ShowDeviceHandler) GenerateBlocks(data any) []slack.Block {
	devices, ok := data.(device.DevicesInfo)
	if !ok {
		log.Fatal("Expected DevicesInfo but got something else")
	}
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
