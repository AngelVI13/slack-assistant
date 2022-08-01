package modals

const (
	MParkingTitle    = "Reserve parking space"
	MParkingActionId = "parkingActionId"
	MParkingOptionId = "parkingOptionId"

	bookParkingAction ModalAction = "book"

	// Default
	DefaultParkingAction ModalAction = bookParkingAction
)

var ParkingActionMap = map[ModalAction]ModalData{
	bookParkingAction: {
		handler:     &ParkingBookingHandler{},
		description: "Reserve a parking space",
	},
}

var ParkingModalInfo = ModalInfo{
	Title:    MParkingTitle,
	ActionId: MParkingActionId,
	OptionId: MParkingOptionId,
}
