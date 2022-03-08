package modals

import (
	"github.com/slack-go/slack"
)

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

// getTakenDevices Get slice of all currently taken devices
func getTakenDevices(devicesInfo DevicesInfo) DevicesInfo {
	var takenDevices DevicesInfo

	for _, device := range devicesInfo {
		if device.Reserved {
			takenDevices = append(takenDevices, device)
		}
	}

	return takenDevices
}

// getAllDevices Get a slice of all devices (copies)
func getAllDevices(devicesInfo DevicesInfo) DevicesInfo {
	allDevices := make(DevicesInfo, len(devicesInfo))
	copy(allDevices, devicesInfo)
	return allDevices
}

// generateDeviceBlocks Generates option block objects to be used as poll elements in modal
func generateDeviceBlocks(devices DevicesInfo, filter func(DevicesInfo) DevicesInfo) []*slack.OptionBlockObject {
	var deviceBlocks []*slack.OptionBlockObject

	for _, device := range filter(devices) {
		sectionBlock := slack.NewOptionBlockObject(
			device.Name,
			slack.NewTextBlockObject("plain_text", device.Name, false, false),
			nil, // TODO: maybe add description
		)
		deviceBlocks = append(deviceBlocks, sectionBlock)
	}

	return deviceBlocks
}
