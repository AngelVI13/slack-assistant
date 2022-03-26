package modals

import (
	"github.com/AngelVI13/slack-assistant/device"
	"github.com/slack-go/slack"
)

const MRestartProxyTitle = "Restart Proxy"
const MRestartProxyActionId = "proxySelected"
const MRestartProxyCheckboxId = "proxyCheckbox"

/*
const MAutoReleaseCheckboxId = "autoReleaseDeviceCheckbox"
const MAutoReleaseActionId = "autoReleaseDeviceSelected"
*/

type RestartProxyHandler struct{}

func (h *RestartProxyHandler) GenerateModalRequest(command *slack.SlashCommand, devices device.DevicesInfo) slack.ModalViewRequest {
	allBlocks := h.GenerateBlocks(devices)
	return generateModalRequest(MRestartProxyTitle, allBlocks)
}

func (h *RestartProxyHandler) GenerateBlocks(devices device.DevicesInfo) []slack.Block {
	deviceOptionBlocks := generateProxyInfoBlocks(devices)
	if len(deviceOptionBlocks) <= 0 {
		return []slack.Block{
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", "Did not find any proxy information.", false, false),
				nil,
				nil,
			),
		}
	}

	// Turn device blocks to a poll/action element block
	deviceCheckboxGroup := slack.NewCheckboxGroupsBlockElement(MRestartProxyCheckboxId, deviceOptionBlocks...)
	actionBlocks := slack.NewActionBlock(MRestartProxyActionId, deviceCheckboxGroup)

	// Header text
	header := "Choose proxy/ies you would like to restart"
	headerText := slack.NewTextBlockObject("mrkdwn", header, false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	/*
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
	*/

	// allBlocks := []slack.Block{headerSection, actionBlocks, divSection, autoReleaseActionBlocks}
	allBlocks := []slack.Block{headerSection, actionBlocks}
	return allBlocks
}
