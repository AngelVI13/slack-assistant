package users

import "github.com/AngelVI13/slack-assistant/device"

type User struct {
	Id         string
	Rights     device.AccessRight
	IsReviewer bool
}

type UsersInfo map[string]*User

type Reviewer struct {
	Name string
	Id   string
}
