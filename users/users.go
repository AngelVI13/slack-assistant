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

	for _, userInfo := range selectedUsersInfo {
		userName := userInfo.Name
		log.Printf("Adding %s", userName)

		u.Map[userName] = &User{
			Id:         userInfo.ID,
			Rights:     accessRights,
			IsReviewer: isReviewer,
		}
	}

	u.SynchronizeToFile()
}

func (u *UsersInfo) IsSpecial(userName string) bool {
	user, ok := u.Map[userName]
	if !ok {
		return false
	}
	return user.Rights == ADMIN
}
