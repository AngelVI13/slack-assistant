package modals

import (
	"github.com/slack-go/slack"
)

const MReserveDeviceTitle = "Reserve Device"

func generateDeviceBlocks() []slack.Block {
	firstNameText := slack.NewTextBlockObject("plain_text", "First Name", false, false)
	firstNamePlaceholder := slack.NewTextBlockObject("plain_text", "Enter your first name", false, false)
	firstNameElement := slack.NewPlainTextInputBlockElement(firstNamePlaceholder, "firstName")
	// Notice that blockID is a unique identifier for a block
	firstName := slack.NewInputBlock("First Name", firstNameText, firstNameElement)

	lastNameText := slack.NewTextBlockObject("plain_text", "Last Name", false, false)
	lastNamePlaceholder := slack.NewTextBlockObject("plain_text", "Enter your first name", false, false)
	lastNameElement := slack.NewPlainTextInputBlockElement(lastNamePlaceholder, "lastName")
	lastName := slack.NewInputBlock("Last Name", lastNameText, lastNameElement)

	return []slack.Block{firstName, lastName}
}

func GenerateModalRequest() slack.ModalViewRequest {
	// Create a ModalViewRequest with a header and two inputs
	titleText := slack.NewTextBlockObject("plain_text", MReserveDeviceTitle, false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Close", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Submit", false, false)

	headerText := slack.NewTextBlockObject("mrkdwn", "Choose a device you would like to reserve", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	deviceBlocks := generateDeviceBlocks()
	allBlocks := []slack.Block{}
	allBlocks = append(allBlocks, headerSection)
	allBlocks = append(allBlocks, deviceBlocks...)

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
