package modals

import (
	"log"

	"github.com/AngelVI13/slack-assistant/device"
	"github.com/slack-go/slack"
)

const MDeviceTitle = "Devices"
const MDeviceActionId = "devicesActionId"
const MDeviceOptionId = "devicesOptionId"

type ModalAction string

const (
	ShowDevicesAction    ModalAction = "show"
	ReserveDevicesAction ModalAction = "reserve"
	ReleaseDevicesAction ModalAction = "release"

	// Default
	DefaultAction ModalAction = ShowDevicesAction
)

var DeviceModalHandlers = map[ModalAction]ModalHandler{
	ShowDevicesAction:    &ShowDeviceHandler{},
	ReserveDevicesAction: &ReserveDeviceHandler{},
	ReleaseDevicesAction: &ReleaseDeviceHandler{},
}

type DeviceHandler struct {
	selectedAction ModalAction
}

func NewDeviceHandler() *DeviceHandler {
	return &DeviceHandler{
		selectedAction: ShowDevicesAction,
	}
}

func (h *DeviceHandler) ChangeAction(action string) {
	modalAction := ModalAction(action)
	_, ok := DeviceModalHandlers[modalAction]
	if !ok {
		log.Fatalf("No such device action exists %s", action)

	}
	h.selectedAction = modalAction
}

func (h *DeviceHandler) Reset() {
	h.selectedAction = DefaultAction
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
		optionText := string(option)
		optionBlock := slack.NewOptionBlockObject(
			optionText,
			slack.NewTextBlockObject("plain_text", optionText, false, false),
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
