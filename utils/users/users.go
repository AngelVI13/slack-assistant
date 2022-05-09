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

type UserMap map[string]AccessRight

func (users *UserMap) SynchronizeToFile() {
	data, err := json.Marshal(users)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(os.Getenv("SL_USERS_FILE"), data, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("INFO: Wrote users list to file")
}
