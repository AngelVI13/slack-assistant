package modals

type ModalAction string

type ModalData struct {
	handler     ModalHandler
	description string
}

type ActionMap map[ModalAction]ModalData

// ModalInfo Info that helps identify which action/option has been changed/selected by user
type ModalInfo struct {
	Title    string
	ActionId string
	OptionId string
}
