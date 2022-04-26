package utils

import (
	"encoding/json"
	"log"
	"os"

	"github.com/AngelVI13/slack-assistant/users"
)

func GetUsers(path string) (users users.UsersInfo) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not read users file (%s)", path)
	}

	err = json.Unmarshal(fileData, &users)
	if err != nil {
		log.Fatalf("Could not parse users file (%s). Error: %+v", path, err)
	}

	loadedUsersNum := len(users)
	if loadedUsersNum == 0 {
		log.Fatalf("No users found in (%s).", path)
	}

	log.Printf("INIT: User list loaded successfully (%d users configured)", loadedUsersNum)
	return users
}
