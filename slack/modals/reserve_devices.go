package modals

import (
	"github.com/AngelVI13/slack-assistant/device"
	"github.com/slack-go/slack"
)

const MReserveDeviceTitle = "Reserve Device"
const MReserveDeviceActionId = "deviceSelected"
const MReserveDeviceCheckboxId = "deviceCheckbox"
const MAutoReleaseCheckboxId = "autoReleaseDeviceCheckbox"
const MAutoReleaseActionId = "autoReleaseDeviceSelected"

type ReserveDeviceHandler struct{}

func (h *ReserveDeviceHandler) GenerateModalRequest(command *slack.SlashCommand, devices device.DevicesInfo) slack.ModalViewRequest {
	allBlocks := h.GenerateBlocks(devices)
	return generateModalRequest(MReserveDeviceTitle, allBlocks)
}

func (h *ReserveDeviceHandler) GenerateBlocks(devices device.DevicesInfo) []slack.Block {
	deviceOptionBlocks := generateDeviceFreeOptionBlocks(devices)
	// If no devices are taken -> return a simple message to the user
	if len(deviceOptionBlocks) <= 0 {
		return []slack.Block{
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", "All devices are taken, nothing to reserve", false, false),
				nil,
				nil,
			),
		}
	}

	// Turn device blocks to a poll/action element block
	deviceCheckboxGroup := slack.NewCheckboxGroupsBlockElement(MReserveDeviceCheckboxId, deviceOptionBlocks...)
	actionBlocks := slack.NewActionBlock(MReserveDeviceActionId, deviceCheckboxGroup)

	// Header text
	header := "Choose a device you would like to reserve"
	headerText := slack.NewTextBlockObject("mrkdwn", header, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Auto release checkbox
	autoReleaseOptionBlocks := []*slack.OptionBlockObject{
		slack.NewOptionBlockObject(
			"autoRelease",
			slack.NewTextBlockObject("mrkdwn", "Auto Release", false, false),
			slack.NewTextBlockObject("mrkdwn", "Automatically release device/s at midnight.", false, false),
		),
	}

	autoReleaseDeviceCheckboxGroup := slack.NewCheckboxGroupsBlockElement(MAutoReleaseCheckboxId, autoReleaseOptionBlocks...)
	autoReleaseActionBlocks := slack.NewActionBlock(MAutoReleaseActionId, autoReleaseDeviceCheckboxGroup)

	divSection := slack.NewDividerBlock()

	allBlocks := []slack.Block{headerSection, actionBlocks, divSection, autoReleaseActionBlocks}
	return allBlocks
}