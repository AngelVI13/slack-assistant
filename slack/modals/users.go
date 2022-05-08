package modals

const (
	MUsersTitle    = "Users"
	MUsersActionId = "UsersActionId"
	MUsersOptionId = "UsersOptionId"

	showUsersAction   ModalAction = "show"
	addUsersAction    ModalAction = "add"
	removeUsersAction ModalAction = "remove"

	// Default
	DefaultUsersAction ModalAction = showUsersAction
)

var UsersActionMap = map[ModalAction]ModalData{
	showUsersAction: {
		handler:     &ShowUsersHandler{},
		description: "Show users",
	},
	addUsersAction: {
		handler:     &AddUserHandler{},
		description: "Add users",
	},
	removeUsersAction: {
		handler:     &RemoveUsersHandler{},
		description: "Remove users",
	},
}

var UsersModalInfo = ModalInfo{
	Title:    MUsersTitle,
	ActionId: MUsersActionId,
	OptionId: MUsersOptionId,
}
