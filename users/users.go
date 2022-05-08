package users

import (
	"math/rand"
	"time"

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
	allReviewers := GetReviewers(usersInfo)

	reviewers := Reviewers{All: allReviewers}
	reviewers.ResetCurrentReviewers()

	return reviewers
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

func removeByIdx(s []*Reviewer, i int) []*Reviewer {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (r *Reviewers) ResetCurrentReviewers() {
	currentReviewers := []*Reviewer{}
	for _, r := range r.All {
		currentReviewers = append(currentReviewers, r)
	}

	r.Current = currentReviewers
}

// ChooseReviewer Picks a reviewer from the current list of reviewers (can't be equal to sender).
// If current reviewer list is empty -> reloads the reviewers list.
func (r *Reviewers) ChooseReviewer(senderName string) *Reviewer {
	rand.Seed(time.Now().UnixNano())

	if len(r.Current) == 0 || (len(r.Current) == 1 && r.Current[0].Name == senderName) {
		r.ResetCurrentReviewers()
	}

	var chosenIdx int
	for {
		chosenIdx = rand.Intn(len(r.Current))

		if r.Current[chosenIdx].Name != senderName {
			break
		}
	}

	reviewer := r.Current[chosenIdx]
	r.Current = removeByIdx(r.Current, chosenIdx)

	return reviewer
}
