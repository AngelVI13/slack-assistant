package modals

import (
	"log"

	"github.com/slack-go/slack"
)

type CustomOptionModalHandler struct {
	selectedAction ModalAction
	defaultAction  ModalAction
	actionMap      ActionMap
	modalInfo      ModalInfo
}

func NewCustomOptionModalHandler(actionMap ActionMap, defaultAction ModalAction, actionInfo ModalInfo) *CustomOptionModalHandler {
	return &CustomOptionModalHandler{
		selectedAction: defaultAction,
		defaultAction:  defaultAction,
		actionMap:      actionMap,
		modalInfo:      actionInfo,
	}
}

func (h *CustomOptionModalHandler) ChangeAction(action string) bool {
	modalAction := ModalAction(action)
	_, ok := h.actionMap[modalAction]
	if !ok {
		return ok
	}
	h.selectedAction = modalAction
	return ok
}

func (h *CustomOptionModalHandler) Reset() {
	h.selectedAction = h.defaultAction
}

func (h *CustomOptionModalHandler) GenerateModalRequest(command *slack.SlashCommand, data any) slack.ModalViewRequest {
	allBlocks := h.GenerateBlocks(command, data)

	var title string
	submition_req := false

	switch h.selectedAction {
	case showUsersAction:
		title = MShowUsersTitle
	case addUsersAction:
		title = MAddUserTitle
		submition_req = true
	case removeUsersAction:
		title = MRemoveUsersTitle
		submition_req = true
	default:
		title = h.modalInfo.Title
	}

	// Input blocks require submit button
	if submition_req {
		return generateModalRequest(title, allBlocks)
	} else {
		return generateInfoModalRequest(title, allBlocks)
	}

}

func (h *CustomOptionModalHandler) GenerateBlocks(command *slack.SlashCommand, data any) []slack.Block {
	var allBlocks []slack.Block

	// Options
	var optionBlocks []*slack.OptionBlockObject

	for option, data := range h.actionMap {
		optionText := string(option)
		optionBlock := slack.NewOptionBlockObject(
			optionText,
			slack.NewTextBlockObject("plain_text", optionText, false, false),
			slack.NewTextBlockObject("plain_text", data.description, false, false),
		)
		optionBlocks = append(optionBlocks, optionBlock)
	}

	// Text shown as title when option box is opened/expanded
	optionLabel := slack.NewTextBlockObject("plain_text", "Action to perform", false, false)
	// Default option shown for option box
	defaultOption := slack.NewTextBlockObject("plain_text", string(h.defaultAction), false, false)

	optionGroupBlockObject := slack.NewOptionGroupBlockElement(optionLabel, optionBlocks...)
	newOptionsGroupSelectBlockElement := slack.NewOptionsGroupSelectBlockElement("static_select", defaultOption, h.modalInfo.OptionId, optionGroupBlockObject)

	actionBlock := slack.NewActionBlock(h.modalInfo.ActionId, newOptionsGroupSelectBlockElement)
	allBlocks = append(allBlocks, actionBlock, slack.NewDividerBlock())

	// Actual modal blocks
	action, ok := h.actionMap[h.selectedAction]
	if !ok {
		log.Fatalf("No such action exists %s", h.selectedAction)
	}
	blocks := action.handler.GenerateBlocks(command, data)
	allBlocks = append(allBlocks, blocks...)

	return allBlocks
}
