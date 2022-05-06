package users

import (
	"log"

	"github.com/AngelVI13/slack-assistant/device"
)

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

type Reviewers struct {
	All     []*Reviewer
	Current []*Reviewer
}

func NewReviewers(usersInfo *UsersInfo) Reviewers {
	// TODO: load this from reviewers file. If doesn't exist create a new file from usersInfo
	reviewers := GetReviewers(usersInfo)
	currentReviewers := make([]*Reviewer, len(reviewers))
	coppiedElems := copy(currentReviewers, reviewers)
	if coppiedElems != len(reviewers) {
		log.Fatalf("failed to copy reviewers to currentReviewers")
	}

	return Reviewers{
		All:     reviewers,
		Current: currentReviewers,
	}
}

func GetReviewers(usersInfo *UsersInfo) (reviewers []*Reviewer) {
	for name, props := range *usersInfo {
		if !props.IsReviewer {
			continue
		}

		reviewers = append(reviewers, &Reviewer{Name: name, Id: props.Id})
	}
	return reviewers
}

func remove(s []*Reviewer, i int) []*Reviewer {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// TODO: add function for choosingReviewer
