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

func getDevicesInfo() []DeviceInfo {
	return []DeviceInfo{
		DeviceInfo{"splinter", false},
		DeviceInfo{"shredder", false},
		DeviceInfo{"donatello", true},
	}
}

// getFreeDevices Get slice of all currently free devices
func getFreeDevices() []DeviceInfo {
	var freeDevices []DeviceInfo

	for _, device := range getDevicesInfo() {
		if !device.Reserved {
			freeDevices = append(freeDevices, device)
		}
	}

	return freeDevices
}

// generateDeviceBlocks Generates option block objects to be used as poll elements in modal
func generateDeviceBlocks() []*slack.OptionBlockObject {
	var deviceBlocks []*slack.OptionBlockObject

	for _, device := range getFreeDevices() {
		sectionBlock := slack.NewOptionBlockObject(
			device.Name,
			slack.NewTextBlockObject("plain_text", device.Name, false, false),
			nil, // TODO: maybe add description
		)
		deviceBlocks = append(deviceBlocks, sectionBlock)
	}

	return deviceBlocks
}

func GenerateModalRequest() slack.ModalViewRequest {
	// Create a ModalViewRequest with a header and two inputs
	titleText := slack.NewTextBlockObject("plain_text", MReserveDeviceTitle, false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Close", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Submit", false, false)

	headerText := slack.NewTextBlockObject("mrkdwn", "Choose a device you would like to reserve", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	deviceOptionBlocks := generateDeviceBlocks()
	// Turn device blocks to a poll/action element block
	deviceCheckboxGroup := slack.NewCheckboxGroupsBlockElement(MReserveDeviceCheckboxId, deviceOptionBlocks...)
	actionBlocks := slack.NewActionBlock(MReserveDeviceActionId, deviceCheckboxGroup)

	// Add header text and action(poll) elem to slice of modal blocks
	allBlocks := []slack.Block{headerSection, actionBlocks}

	blocks := slack.Blocks{
		BlockSet: allBlocks,
	}

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	return modalRequest
}
