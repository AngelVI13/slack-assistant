package users

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/AngelVI13/slack-assistant/config"
)

type AccessRight int

// NOTE: Currently access rights are not used
const (
	STANDARD AccessRight = iota
	ADMIN
)

type User struct {
	Id         string
	Rights     AccessRight
	IsReviewer bool `json:"is_reviewer"`
}

type UsersMap map[string]*User

type UsersInfo struct {
	Map      UsersMap
	Filename string
}

func (u *UsersInfo) SynchronizeToFile() {
	data, err := json.Marshal(u)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(u.Filename, data, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("INFO: Wrote users list to file")
}

type Reviewer struct {
	Name string
	Id   string
}

// TODO: Update this data when add/edit/remove users feature is done
type Reviewers struct {
	All       []*Reviewer
	Current   []*Reviewer
	Filename  string
	ChannelId string // where to post chosen reviewer messages
}

func NewReviewers(config *config.Config, usersMap *UsersMap) Reviewers {
	filename := config.ReviewersFilename

	allReviewers := GetReviewers(usersMap)
	reviewers := Reviewers{All: allReviewers, Filename: filename, ChannelId: config.SlackTaChannelId}

	_, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		// Load reviewers from users list and create current reviewers list file
		reviewers.ResetCurrentReviewers()
		reviewers.synchronizeToFile()

		log.Printf("INFO: Generated reviewers list from users info (%d reviewers).", len(reviewers.All))
	} else if err == nil {
		// Load current reviewers from file
		// All reviewers info comes from users info
		data, err := os.ReadFile(filename)
		if err != nil {
			log.Fatalf("Failed to read from reviewers file: %+v", err)
		}

		reviewers.synchronizeFromFile(data)
	} else {
		log.Fatalf("Initializing reviewers failed. Couldn't open current reviewers file: %+v", err)
	}

	return reviewers
}

func GetReviewers(usersInfo *UsersMap) (reviewers []*Reviewer) {
	for name, props := range *usersInfo {
		if !props.IsReviewer {
			continue
		}

		reviewers = append(reviewers, &Reviewer{Name: name, Id: props.Id})
	}
	return reviewers
}

// TODO: Make this into a generic if its used more often
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

	r.synchronizeToFile()

	return reviewer
}

func (r *Reviewers) synchronizeToFile() {
	data, err := json.Marshal(r.Current)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(r.Filename, data, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("INFO: Wrote reviewers list to file")
}

func (r *Reviewers) synchronizeFromFile(data []byte) {
	err := json.Unmarshal(data, &r.Current)
	if err != nil {
		log.Fatalf("Could not parse reviewers file. Error: %+v", err)
	}
	log.Printf("INFO: Reviewers list loaded successfully (%d reviewers in queue)", len(r.Current))
	log.Printf("INFO: ---------------------------------- (%d total reviewers)", len(r.All))
}
