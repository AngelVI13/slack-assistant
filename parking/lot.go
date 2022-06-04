package parking

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/AngelVI13/slack-assistant/config"
)

type SpacesInfo []*ParkingSpace

type ParkingLot struct {
	ParkingSpaces
	config *config.Config
}

func NewParkingLot() ParkingLot {
	return ParkingLot{
		ParkingSpaces: make(map[int]*ParkingSpace),
	}
}

// NewParkingLotFromJson Takes json data as input and returns a populated ParkingLot object
func NewParkingLotFromJson(data []byte, config *config.Config) ParkingLot {
	parkingLot := NewParkingLot()
	parkingLot.synchronizeFromFile(data)
	parkingLot.config = config
	return parkingLot
}

func (d *ParkingLot) SynchronizeToFile() {
	data, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(d.config.ParkingFilename, data, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("INFO: Wrote parking lot to file")
}

func (d *ParkingLot) synchronizeFromFile(data []byte) {
	// Unmarshal the provided data into the solid map
	err := json.Unmarshal(data, d)
	if err != nil {
		log.Fatalf("Could not parse parking file. Error: %+v", err)
	}
}

// TODO: This is identical to GetDevicesInfo -> refactor it out
func (d *ParkingLot) GetSpacesInfo(user string) SpacesInfo {
	// Group devices in 2 groups -> belonging to given user or not
	// The group that doesn't belong to user will be sorted by name and by status (reserved or not)
	userSpaces := make(SpacesInfo, 0)
	nonUserSpaces := make(SpacesInfo, 0)
	for _, d := range d.ParkingSpaces {
		if d.Reserved && d.ReservedBy == user {
			userSpaces = append(userSpaces, d)
		} else {
			nonUserSpaces = append(nonUserSpaces, d)
		}
	}

	// NOTE: This sorts the device list starting from free devices
	sort.Slice(nonUserSpaces, func(i, j int) bool {
		return !nonUserSpaces[i].Reserved
	})

	firstTaken := -1 // Index of first taken device
	for i, device := range nonUserSpaces {
		if device.Reserved {
			firstTaken = i
			break
		}
	}

	// NOTE: this might be unnecessary but it shows devices in predicable way in UI so its nice.
	// If all devices are free or all devices are taken, sort by name
	if firstTaken == -1 || firstTaken == 0 {
		sort.Slice(nonUserSpaces, func(i, j int) bool {
			return nonUserSpaces[i].Number < nonUserSpaces[j].Number
		})
	} else {
		// split devices into 2 - free & taken
		// sort each sub slice based on device name/port
		free := nonUserSpaces[:firstTaken]
		taken := nonUserSpaces[firstTaken:]

		sort.Slice(free, func(i, j int) bool {
			return free[i].Number < free[j].Number
		})

		sort.Slice(taken, func(i, j int) bool {
			return taken[i].Number < taken[j].Number
		})
	}

	allSpaces := make(SpacesInfo, 0, len(d.ParkingSpaces))
	allSpaces = append(allSpaces, userSpaces...)
	allSpaces = append(allSpaces, nonUserSpaces...)
	return allSpaces
}

func (l *ParkingLot) Reserve(parkingSpace, user, userId string, autoRelease bool) (errMsg string) {
	spaceNumber, err := strconv.Atoi(parkingSpace)
	if err != nil {
		log.Fatalf("Could not convert parkingSpace %+v to int", parkingSpace)
	}

	space, ok := l.ParkingSpaces[spaceNumber]
	if !ok {
		log.Fatalf("Wrong parking space number %d, %+v", spaceNumber, l)
	}
	// Only inform user if it was someone else that tried to reserved his device.
	// This prevents an unnecessary message if you double clicked the reserve button yourself
	if space.Reserved && space.ReservedById != userId {
		reservedTime := space.ReservedTime.Format("Mon 15:04")
		return fmt.Sprintf("*Error*: Could not reserve *%d*. *%s* has just reserved it (at *%s*)", spaceNumber, space.ReservedBy, reservedTime)
	}
	log.Printf("PARKING_RESERVE: User (%s) reserved space (%d) with auto release (%v)", user, spaceNumber, autoRelease)

	space.Reserved = true
	space.ReservedBy = user
	space.ReservedById = userId
	space.ReservedTime = time.Now()
	space.AutoRelease = autoRelease

	l.SynchronizeToFile()
	return ""
}

func (l *ParkingLot) Release(parkingSpace, user string) (victimId, errMsg string) {
	spaceNumber, err := strconv.Atoi(parkingSpace)
	if err != nil {
		log.Fatalf("Could not convert parkingSpace %+v to int", parkingSpace)
	}

	space, ok := l.ParkingSpaces[spaceNumber]
	if !ok {
		log.Fatalf("Incorrect parking space number %s, %+v", parkingSpace, l)
	}

	log.Printf("PARKING_RELEASE: User (%s) released (%s) device.", user, parkingSpace)

	space.Reserved = false
	l.SynchronizeToFile()

	if space.ReservedBy != user {
		return space.ReservedById, fmt.Sprintf(":warning: *%s* released your (*%s*) space (*%d*)", user, space.ReservedBy, space.Number)
	}
	return "", ""
}

func (l *ParkingLot) AutoRelease(when time.Time) {
	// Only release devices at the specified hour (hour is [0;23])
	now := time.Now()
	if now.Hour() != when.Hour() {
		return
	}

	for _, space := range l.ParkingSpaces {
		if space.Reserved && space.AutoRelease {
			space.Reserved = false
			space.AutoRelease = false
		}
	}

	// Need to synchronize changes from file otherwise the state won't be preserved after restart
	// NOTE: This ends up synchronizing to file more than once since the function can be called
	// multiple times within the specified auto release hour (even if nothing has changed in the devices list).
	l.SynchronizeToFile()
}
