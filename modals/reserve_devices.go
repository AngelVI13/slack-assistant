package modals

import (
	"github.com/slack-go/slack"
)

const MReserveDeviceTitle = "Reserve Device"
const MReserveDeviceActionId = "deviceSelected"
const MReserveDeviceCheckboxId = "deviceCheckbox"

type DeviceInfo struct {
	Name     string
	Reserved bool
}

type DevicesInfo []DeviceInfo

// Devices Global variable holding all device info
var Devices DevicesInfo = DevicesInfo{
	DeviceInfo{"splinter", false},
	DeviceInfo{"shredder", false},
	DeviceInfo{"donatello", true},
}

// getFreeDevices Get slice of all currently free devices
func getFreeDevices(devicesInfo DevicesInfo) DevicesInfo {
	var freeDevices DevicesInfo

	for _, device := range devicesInfo {
		if !device.Reserved {
			freeDevices = append(freeDevices, device)
		}
	}

	return freeDevices
}

// generateDeviceBlocks Generates option block objects to be used as poll elements in modal
func generateDeviceBlocks(devices DevicesInfo) []*slack.OptionBlockObject {
	var deviceBlocks []*slack.OptionBlockObject

	for _, device := range getFreeDevices(devices) {
		sectionBlock := slack.NewOptionBlockObject(
			device.Name,
			slack.NewTextBlockObject("plain_text", device.Name, false, false),
			nil, // TODO: maybe add description
		)
		deviceBlocks = append(deviceBlocks, sectionBlock)
	}

	return deviceBlocks
}

func GenerateReserveDeviceModalRequest(devices DevicesInfo) slack.ModalViewRequest {
	deviceOptionBlocks := generateDeviceBlocks(devices)
	// Turn device blocks to a poll/action element block
	deviceCheckboxGroup := slack.NewCheckboxGroupsBlockElement(MReserveDeviceCheckboxId, deviceOptionBlocks...)
	actionBlocks := slack.NewActionBlock(MReserveDeviceActionId, deviceCheckboxGroup)

	header := "Choose a device you would like to reserve"
	headerText := slack.NewTextBlockObject("mrkdwn", header, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Add header text and action(poll) elem to slice of modal blocks
	allBlocks := []slack.Block{headerSection, actionBlocks}
	return generateModalRequest(MReserveDeviceTitle, allBlocks)
}
