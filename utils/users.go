package utils

import (
	"encoding/json"
	"github.com/AngelVI13/slack-assistant/device"
	"log"
	"os"
)

func GetUsers(path string) (usersList map[string]device.AccessRight) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not read users file (%s)", path)
	}

	err = json.Unmarshal(fileData, &usersList)
	if err != nil {
		log.Fatalf("Could not parse users file (%s). Error: %+v", path, err)
	}

	loadedUsersNum := len(usersList)
	if loadedUsersNum == 0 {
		log.Fatalf("No users found in (%s).", path)
	}

	log.Printf("INIT: User list loaded successfully (%d users configured)", loadedUsersNum)
	return usersList
}
