package modals

import (
	"fmt"
	"log"

	"github.com/AngelVI13/slack-assistant/device"
	"github.com/slack-go/slack"
)

const MRestartProxyTitle = "Restart Proxy"
const MRestartProxyActionId = "proxySelected"
const MRestartProxyCheckboxId = "proxyCheckbox"

type RestartProxyHandler struct{}

func (h *RestartProxyHandler) GenerateModalRequest(command *slack.SlashCommand, data ...any) slack.ModalViewRequest {
	allBlocks := h.GenerateBlocks(command, data...)
	return generateModalRequest(MRestartProxyTitle, allBlocks)
}

func (h *RestartProxyHandler) GenerateBlocks(command *slack.SlashCommand, data ...any) []slack.Block {
	if len(data) != 1 {
		log.Fatal("Incorrect num of params for RestartProxyHandler: 1")
	}
	rawDevices := data[0]

	devices, ok := rawDevices.(device.DevicesInfo)
	if !ok {
		log.Fatal("Expected DevicesInfo but got something else")
	}

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

	allBlocks := []slack.Block{headerSection, actionBlocks}
	return allBlocks
}

// generateProxyInfoBlocks Generates device proxy block objects to be used as elements in modal
func generateProxyInfoBlocks(devices device.DevicesInfo) []*slack.OptionBlockObject {
	var deviceBlocks []*slack.OptionBlockObject

	for _, device := range devices {
		status := device.GetStatusDescription()
		emoji := device.GetStatusEmoji()

		deviceProps := device.GetPropsText()
		text := fmt.Sprintf("%s\n%s", deviceProps, status)
		optionName := fmt.Sprintf("%s %s", emoji, device.Name)

		sectionBlock := slack.NewOptionBlockObject(
			device.Name,
			slack.NewTextBlockObject("mrkdwn", optionName, false, false),
			slack.NewTextBlockObject("mrkdwn", text, false, false),
		)
		deviceBlocks = append(deviceBlocks, sectionBlock)
	}

	return deviceBlocks
}
