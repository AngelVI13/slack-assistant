package users

import (
	"encoding/json"
	"log"
	"os"
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
	data, err := json.Marshal(u.Map)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(u.Filename, data, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("INFO: Wrote users list to file")
}
