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

func (h *ShowDeviceHandler) GenerateModalRequest(data any) slack.ModalViewRequest {
	allBlocks := h.GenerateBlocks(data)
	return generateInfoModalRequest(MShowDeviceTitle, allBlocks)
}

func (h *ShowDeviceHandler) GenerateBlocks(data any) []slack.Block {
	devices, ok := data.(device.DevicesInfo)
	if !ok {
		log.Fatal("Expected DevicesInfo but got something else")
	}
	deviceSectionBlocks := generateDeviceInfoBlocks(devices)

	// TODO: 1. remove Reserve and Release device modals cause now everything can be done via ShowDeviceModal
	// 2. which should also be renamed to Devices or sth
	// 3. Also remove the devices from the OptionModal
	// 4. Also connect the button logic to do something
	allBlocks := deviceSectionBlocks
	return allBlocks
}
