package users

import (
	"encoding/json"
	"log"
	"os"

	"github.com/slack-go/slack"
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
	data, err := json.MarshalIndent(u.Map, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(u.Filename, data, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("INFO: Wrote users list to file")
}

func (u *UsersInfo) AddNewUsers(selectedUsersInfo []*slack.User, selectedOptions []slack.OptionBlockObject, accsessRightSelection string, reviewerOptionSelection string) {
	accessRights := STANDARD
	isReviewer := false

	for _, selection := range selectedOptions {
		switch selection.Value {
		case accsessRightSelection:
			accessRights = ADMIN
		case reviewerOptionSelection:
			isReviewer = true
		}
	}

	for _, user_info := range selectedUsersInfo {
		user_name := user_info.Name
		log.Printf("Adding %s", user_name)

		u.Map[user_name] = &User{
			Id:         user_info.ID,
			Rights:     accessRights,
			IsReviewer: isReviewer,
		}
	}

	u.SynchronizeToFile()

}
