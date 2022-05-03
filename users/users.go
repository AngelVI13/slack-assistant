package users

import "github.com/AngelVI13/slack-assistant/device"

type User struct {
	Id         string
	Rights     device.AccessRight
	IsReviewer bool `json:"is_reviewer"`
}

type UsersInfo map[string]*User

type Reviewer struct {
	Name string
	Id   string
}

type Reviewers []*Reviewer
