package modals

import (
	"log"

	"github.com/AngelVI13/slack-assistant/device"
	"github.com/slack-go/slack"
)

const MRestartProxyTitle = "Restart Proxy"
const MRestartProxyActionId = "proxySelected"
const MRestartProxyCheckboxId = "proxyCheckbox"

type RestartProxyHandler struct{}

func (h *RestartProxyHandler) GenerateModalRequest(data any) slack.ModalViewRequest {
	allBlocks := h.GenerateBlocks(data)
	return generateModalRequest(MRestartProxyTitle, allBlocks)
}

func (h *RestartProxyHandler) GenerateBlocks(data any) []slack.Block {
	devices, ok := data.(device.DevicesInfo)
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
