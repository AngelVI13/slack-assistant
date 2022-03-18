package modals

import (
	"fmt"
	"github.com/slack-go/slack"
	"time"
)

type WorkerProps struct {
	SerialProxyHost    string `json:"serial_proxy_host"`
	SerialProxyPort    int    `json:"serial_proxy_port"`
	DeviceSerialNumber string `json:"device_serial_number"`
}

type DeviceProps struct {
	Name         string
	Reserved     bool
	ReservedBy   string
	ReservedById string
	ReservedTime time.Time
	Props        WorkerProps
}

type DevicesInfo []*DeviceProps

func (p *DeviceProps) GetPropsText() string {
	return fmt.Sprintf(
		"Port: %d\tHost: %s\tS/N: %s",
		p.Props.SerialProxyPort,
		p.Props.SerialProxyHost,
		p.Props.DeviceSerialNumber,
	)
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

// generateDeviceFreeOptionBlocks Generates option block objects for every free device
// to be used as poll elements in modal
func generateDeviceFreeOptionBlocks(devices DevicesInfo) []*slack.OptionBlockObject {
	var deviceBlocks []*slack.OptionBlockObject

	for _, device := range getFreeDevices(devices) {
		description := device.GetPropsText()
		sectionBlock := slack.NewOptionBlockObject(
			device.Name,
			slack.NewTextBlockObject("mrkdwn", device.Name, false, false),
			slack.NewTextBlockObject("mrkdwn", description, false, false),
		)
		deviceBlocks = append(deviceBlocks, sectionBlock)
	}

	return deviceBlocks
}

// generateDeviceTakenOptionBlocks Generates option block objects for every taken device
// to be used as poll elements in modal
func generateDeviceTakenOptionBlocks(devices DevicesInfo) []*slack.OptionBlockObject {
	var deviceBlocks []*slack.OptionBlockObject

	for _, device := range getTakenDevices(devices) {
		timeStr := device.ReservedTime.Format("Mon 15:04")
		deviceProps := device.GetPropsText()
		status := fmt.Sprintf("%s\nReserved by\t:bust_in_silhouette:*%s*\ton\t:clock1: *%s*", deviceProps, device.ReservedBy, timeStr)

		sectionBlock := slack.NewOptionBlockObject(
			device.Name,
			slack.NewTextBlockObject("mrkdwn", device.Name, false, false),
			slack.NewTextBlockObject("mrkdwn", status, false, false),
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
		emoji := ":large_green_circle:"
		if device.Reserved {
			emoji = ":large_orange_circle:"
			timeStr := device.ReservedTime.Format("Mon 15:04")
			status = fmt.Sprintf("Reserved by\t:bust_in_silhouette:*%s*\ton\t:clock1: *%s*", device.ReservedBy, timeStr)
		}
		deviceProps := device.GetPropsText()
		text := fmt.Sprintf("%s *%s*\n\t\t%s\n\t\t_%s_", emoji, device.Name, deviceProps, status)
		sectionText := slack.NewTextBlockObject("mrkdwn", text, false, false)
		sectionBlock := slack.NewSectionBlock(sectionText, nil, nil)

		deviceBlocks = append(deviceBlocks, sectionBlock)
	}

	return deviceBlocks
}
