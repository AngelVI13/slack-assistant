package users

import "github.com/AngelVI13/slack-assistant/device"

type User struct {
	Name       string
	Id         string
	Rights     device.AccessRight
	IsReviewer bool
}

type UsersInfo map[string]*User
