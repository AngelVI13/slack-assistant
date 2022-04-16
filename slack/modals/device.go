package modals

import (
	"log"

	"github.com/AngelVI13/slack-assistant/device"
	"github.com/slack-go/slack"
)

const (
	MDeviceTitle    = "Devices"
	MDeviceActionId = "devicesActionId"
	MDeviceOptionId = "devicesOptionId"
)

const (
	showDevicesAction    ModalAction = "show"
	reserveDevicesAction ModalAction = "reserve"
	releaseDevicesAction ModalAction = "release"

	// Default
	defaultAction ModalAction = showDevicesAction
)

var deviceModalHandlers = map[ModalAction]ModalData{
	showDevicesAction: {
		handler:     &ShowDeviceHandler{},
		description: "Show list of devices with their status",
	},
	reserveDevicesAction: {
		handler:     &ReserveDeviceHandler{},
		description: "Reserve available devices",
	},
	releaseDevicesAction: {
		handler:     &ReleaseDeviceHandler{},
		description: "Release taken devices",
	},
}

type DeviceHandler struct {
	selectedAction ModalAction
}

func NewDeviceHandler() *DeviceHandler {
	return &DeviceHandler{
		selectedAction: showDevicesAction,
	}
}

func (h *DeviceHandler) ChangeAction(action string) {
	modalAction := ModalAction(action)
	_, ok := deviceModalHandlers[modalAction]
	if !ok {
		log.Fatalf("No such device action exists %s", action)

	}
	h.selectedAction = modalAction
}

func (h *DeviceHandler) Reset() {
	h.selectedAction = defaultAction
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

	for option, data := range deviceModalHandlers {
		optionText := string(option)
		optionBlock := slack.NewOptionBlockObject(
			optionText,
			slack.NewTextBlockObject("plain_text", optionText, false, false),
			slack.NewTextBlockObject("plain_text", data.description, false, false),
		)
		optionBlocks = append(optionBlocks, optionBlock)
	}

	label1 := slack.NewTextBlockObject("plain_text", "Select", false, false)
	placeholder1 := slack.NewTextBlockObject("plain_text", "Select", false, false)
	optionGroupBlockObject := slack.NewOptionGroupBlockElement(label1, optionBlocks...)
	newOptionsGroupSelectBlockElement := slack.NewOptionsGroupSelectBlockElement("static_select", placeholder1, MDeviceOptionId, optionGroupBlockObject)

	actionBlock := slack.NewActionBlock(MDeviceActionId, newOptionsGroupSelectBlockElement)
	allBlocks = append(allBlocks, actionBlock)

	action, ok := deviceModalHandlers[h.selectedAction]
	if !ok {
		log.Fatalf("No such device action exists %s", h.selectedAction)
	}
	blocks := action.handler.GenerateBlocks(devices)
	allBlocks = append(allBlocks, blocks...)

	return allBlocks
}
