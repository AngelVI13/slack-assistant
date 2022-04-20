package modals

const (
	MDeviceTitle    = "Devices"
	MDeviceActionId = "devicesActionId"
	MDeviceOptionId = "devicesOptionId"

	bookDevicesAction    ModalAction = "book"
	restartdevicesAction ModalAction = "restart"

	// Default
	DefaultDeviceAction ModalAction = bookDevicesAction
)

var DeviceActionMap = map[ModalAction]ModalData{
	bookDevicesAction: {
		handler:     &DeviceBookingHandler{},
		description: "Reserve/Release devices",
	},
	/* //TODO: delete this
	restartdevicesAction: {
		handler:     &RestartProxyHandler{},
		description: "Restart proxy for selected devices",
	},
	*/
}

var DeviceModalInfo = ModalInfo{
	Title:    MDeviceTitle,
	ActionId: MDeviceActionId,
	OptionId: MDeviceOptionId,
}
