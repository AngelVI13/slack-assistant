package modals

import (
	"fmt"
	"log"

	"github.com/AngelVI13/slack-assistant/device"
	"github.com/slack-go/slack"
)

const (
	MDeviceBookingTitle = "Device Booking"

	ReserveDeviceActionId   = "reserve"
	ReserveWithAutoActionId = "withAuto"
	ReleaseDeviceActionId   = "release"
)

type DeviceBookingHandler struct{}

func (h *DeviceBookingHandler) GenerateModalRequest(data any) slack.ModalViewRequest {
	allBlocks := h.GenerateBlocks(data)
	// TODO: change modal button to OK instead of submit
	return generateInfoModalRequest(MDeviceBookingTitle, allBlocks)
}

func (h *DeviceBookingHandler) GenerateBlocks(data any) []slack.Block {
	devices, ok := data.(device.DevicesInfo)
	if !ok {
		log.Fatal("Expected DevicesInfo but got something else")
	}
	deviceSectionBlocks := generateDeviceInfoBlocks(devices)

	allBlocks := deviceSectionBlocks
	return allBlocks
}

// generateDeviceInfo Generate sections of text that contain device info such as status (taken/free), ip, port, taken by etc..
func generateDeviceInfo(devices device.DevicesInfo) []slack.SectionBlock {
	var sections []slack.SectionBlock
	for _, device := range devices {
		status := device.GetStatusDescription()
		emoji := device.GetStatusEmoji()

		deviceProps := device.GetPropsText()
		text := fmt.Sprintf("%s *%s*\n\t\t%s\n\t\t%s", emoji, device.Name, deviceProps, status)
		sectionText := slack.NewTextBlockObject("mrkdwn", text, false, false)
		sectionBlock := slack.NewSectionBlock(sectionText, nil, nil)

		sections = append(sections, *sectionBlock)
	}
	return sections
}

func generateDeviceButtons(device *device.DeviceProps) []slack.BlockElement {
	var buttons []slack.BlockElement

	if device.Reserved {
		releaseButton := slack.NewButtonBlockElement(
			ReleaseDeviceActionId,
			device.Name,
			slack.NewTextBlockObject("plain_text", "Release!", true, false),
		)
		releaseButton = releaseButton.WithStyle(slack.StyleDanger)
		buttons = append(buttons, releaseButton)
	} else {

		actionButtonText := "Reserve!"
		reserveWithAutoButton := slack.NewButtonBlockElement(
			ReserveWithAutoActionId,
			device.Name,
			slack.NewTextBlockObject("plain_text", fmt.Sprintf("%s :eject:", actionButtonText), true, false),
		)
		reserveWithAutoButton = reserveWithAutoButton.WithStyle(slack.StylePrimary)
		buttons = append(buttons, reserveWithAutoButton)

		reserveButton := slack.NewButtonBlockElement(ReserveDeviceActionId, device.Name, slack.NewTextBlockObject("plain_text", actionButtonText, true, false))
		buttons = append(buttons, reserveButton)
	}
	return buttons
}

// generateDeviceInfoBlocks Generates device block objects to be used as elements in modal
func generateDeviceInfoBlocks(devices device.DevicesInfo) []slack.Block {
	div := slack.NewDividerBlock()
	deviceSections := generateDeviceInfo(devices)

	var deviceBlocks []slack.Block
	for idx, device := range devices {
		sectionBlock := deviceSections[idx]
		buttons := generateDeviceButtons(device)

		actions := slack.NewActionBlock("", buttons...)
		deviceBlocks = append(deviceBlocks, sectionBlock, actions, div)
	}

	return deviceBlocks
}
