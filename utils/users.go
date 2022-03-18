package utils

import (
	"encoding/json"
	"github.com/AngelVI13/slack-assistant/handlers"
	"log"
	"os"
)

func GetUsers(path string) (usersList map[string]handlers.AccessRight) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not read users file %s", path)
	}

	err = json.Unmarshal(fileData, &usersList)
	if err != nil {
		log.Fatalf("Could not parse users file %s. Error: %+v", path, err)
	}

	log.Printf("User list loaded successfully (%d users configured)", len(usersList))
	return usersList
}
