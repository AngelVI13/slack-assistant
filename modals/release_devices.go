package modals

import (
	"github.com/slack-go/slack"
)

const MReleaseDeviceTitle = "Release Device"
const MReleaseDeviceActionId = "deviceSelected"
const MReleaseDeviceCheckboxId = "deviceCheckbox"

func GenerateReleaseDeviceModalRequest(devices DevicesInfo) slack.ModalViewRequest {
	deviceOptionBlocks := generateDeviceBlocks(devices, getTakenDevices)
	// Turn device blocks to a poll/action element block
	deviceCheckboxGroup := slack.NewCheckboxGroupsBlockElement(MReleaseDeviceCheckboxId, deviceOptionBlocks...)
	actionBlocks := slack.NewActionBlock(MReleaseDeviceActionId, deviceCheckboxGroup)

	header := "Choose a device you would like to release"
	headerText := slack.NewTextBlockObject("mrkdwn", header, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Add header text and action(poll) elem to slice of modal blocks
	allBlocks := []slack.Block{headerSection, actionBlocks}
	return generateModalRequest(MReleaseDeviceTitle, allBlocks)
}
