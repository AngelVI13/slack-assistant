package utils

import (
	"log"
	"os"

	"github.com/AngelVI13/slack-assistant/config"
	"github.com/AngelVI13/slack-assistant/parking"
)

func GetParkingLot(config *config.Config) (parkingLot parking.ParkingLot) {
	path := config.ParkingFilename

	fileData, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not read parking file (%s)", path)
	}

	parkingLot = parking.NewParkingLotFromJson(fileData, config)

	loadedSpacesNum := len(parkingLot.ParkingSpaces)
	if loadedSpacesNum == 0 {
		log.Fatalf("No spaces found in (%s).", path)
	}

	log.Printf("INIT: Parking spaces list loaded successfully (%d spaces configured)", loadedSpacesNum)
	return parkingLot
}
