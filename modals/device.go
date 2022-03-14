package modals

import (
	"fmt"
	"github.com/slack-go/slack"
	"time"
)

type DeviceInfo struct {
	Name         string
	Reserved     bool
	ReservedBy   string
	ReservedTime time.Time
}

type DevicesInfo []*DeviceInfo

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

// generateDeviceOptionBlocks Generates option block objects to be used as poll elements in modal
func generateDeviceOptionBlocks(devices DevicesInfo, filter func(DevicesInfo) DevicesInfo) []*slack.OptionBlockObject {
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

// generateDeviceInfoBlocks Generates device block objects to be used as elements in modal
func generateDeviceInfoBlocks(devices DevicesInfo) []*slack.SectionBlock {
	var deviceBlocks []*slack.SectionBlock

	for _, device := range devices {
		status := "Free"
		emoji := ":large_orange_circle:"
		if device.Reserved {
			emoji = ":large_green_circle:"
			timeStr := device.ReservedTime.Format("Mon 15:04:05")
			status = fmt.Sprintf("Reserved by\t:bust_in_silhouette:*%s*\ton\t:clock1: *%s*", device.ReservedBy, timeStr)
		}
		text := fmt.Sprintf("%s *%s*\n\t\t%s", emoji, device.Name, status)
		sectionText := slack.NewTextBlockObject("mrkdwn", text, false, false)
		sectionBlock := slack.NewSectionBlock(sectionText, nil, nil)

		deviceBlocks = append(deviceBlocks, sectionBlock)
	}

	return deviceBlocks
}
