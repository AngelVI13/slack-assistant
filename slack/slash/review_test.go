package slash

import (
	"os"
	"testing"

	"github.com/AngelVI13/slack-assistant/users"
)

func isInReviewers(reviewers []*users.Reviewer, name string) bool {
	for _, r := range reviewers {
		if r.Name == name {
			return true
		}
	}
	return false
}

func TestChooseReviewer(t *testing.T) {
	// Case1: No reviewers file
	reviewersFile := "does_not_exist.txt"
	reviewersInfo := users.Reviewers{
		All: []*users.Reviewer{
			{Name: "angel.iliev", Id: "1"},
			{Name: "laima.strigo", Id: "2"},
			{Name: "aurimas.razmis", Id: "3"},
		},
		Current: []*users.Reviewer{
			{Name: "angel.iliev", Id: "1"},
			{Name: "laima.strigo", Id: "2"},
			{Name: "aurimas.razmis", Id: "3"},
		},
		Filename: reviewersFile,
	}

	sender := "angel.iliev"

	reviewer := reviewersInfo.ChooseReviewer(sender)
	if reviewer.Name == sender {
		t.Errorf("Chosen reviewer is the same as sender! reviewer == sender == %s", reviewer)
	}

	if !isInReviewers(reviewersInfo.All, reviewer.Name) {
		t.Errorf("Chosen reviewer is not in ALL reviewers list %s; %+v", reviewer.Name, reviewersInfo.All)
	}

	if _, err := os.Stat(reviewersInfo.Filename); err != nil {
		t.Errorf("Reviewers file was not created/updated: %s", reviewersInfo.Filename)
	}

	if err := os.Remove(reviewersInfo.Filename); err != nil {
		t.Errorf("Failed to remove reviewers file after test was finished: %s", reviewersInfo.Filename)
	}
}
