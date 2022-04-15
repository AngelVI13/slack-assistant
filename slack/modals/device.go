package modals

import (
	"log"

	"github.com/AngelVI13/slack-assistant/device"
	"github.com/slack-go/slack"
)

const MDeviceTitle = "Devices"
const MDeviceActionId = "devicesActionId"
const MDeviceOptionId = "devicesOptionId"

var DeviceModalHandlers = map[string]ModalHandler{
	"show":    &ShowDeviceHandler{},
	"reserve": &ReserveDeviceHandler{},
	"release": &ReleaseDeviceHandler{},
}

type DeviceHandler struct {
	selectedAction string
}

func NewDeviceHandler() *DeviceHandler {
	return &DeviceHandler{
		selectedAction: "show",
	}
}

func (h *DeviceHandler) ChangeAction(action string) {
	_, ok := DeviceModalHandlers[action]
	if !ok {
		log.Fatalf("No such device action exists %s", action)

	}
	h.selectedAction = action
}

func (h *DeviceHandler) GenerateModalRequest(data any) slack.ModalViewRequest {
	allBlocks := h.GenerateBlocks(data)
	return generateInfoModalRequest(MDeviceTitle, allBlocks)
}

func (h *DeviceHandler) GenerateBlocks(data any) []slack.Block {
	devices, ok := data.(device.DevicesInfo)
	if !ok {
		log.Fatal("Expected DevicesInfo but got something else")
	}

	var allBlocks []slack.Block

	// Options
	var optionBlocks []*slack.OptionBlockObject

	for option := range DeviceModalHandlers {
		optionBlock := slack.NewOptionBlockObject(
			option,
			slack.NewTextBlockObject("plain_text", option, false, false),
			slack.NewTextBlockObject("plain_text", "description1", false, false),
		)
		optionBlocks = append(optionBlocks, optionBlock)
	}

	label1 := slack.NewTextBlockObject("plain_text", "Select", false, false)
	placeholder1 := slack.NewTextBlockObject("plain_text", "Select", false, false)
	optionGroupBlockObject := slack.NewOptionGroupBlockElement(label1, optionBlocks...)
	newOptionsGroupSelectBlockElement := slack.NewOptionsGroupSelectBlockElement("static_select", placeholder1, MDeviceOptionId, optionGroupBlockObject)

	actionBlock := slack.NewActionBlock(MDeviceActionId, newOptionsGroupSelectBlockElement)
	allBlocks = append(allBlocks, actionBlock)

	actionHandler, ok := DeviceModalHandlers[h.selectedAction]
	if !ok {
		log.Fatalf("No such device action exists %s", h.selectedAction)
	}
	blocks := actionHandler.GenerateBlocks(devices)
	allBlocks = append(allBlocks, blocks...)

	return allBlocks
}
