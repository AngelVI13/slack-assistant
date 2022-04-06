package users

type AccessRight int

// NOTE: Currently access rights are not used
const (
	STANDARD AccessRight = iota
	ADMIN
)

type UserMap map[string]AccessRight
