package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AngelVI13/slack-assistant/handlers"
	"github.com/AngelVI13/slack-assistant/modals"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Worker struct {
	Name        string             `json:"name"`
	Password    string             `json:"password"`
	WorkerSetup []string           `json:"worker_setup"`
	Properties  modals.WorkerProps `json:"properties"`
}

type WorkersResponse struct {
	Workers      []Worker `json:"workers"`
	DeviceSetups []string `json:"device_setups"`
}

func GetTmtWorkers(endpoint string) (*WorkersResponse, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var workers WorkersResponse
	err = json.Unmarshal(body, &workers)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshall workers response %s. Error: %+v", string(body), err)
	}

	log.Printf("Workers info fetched successfully")
	return &workers, nil
}

func GetDevices(path, taProjectEndpoint string) handlers.DevicesMap {
	// 1.a If device info file exists -> read info from there
	// 1.b Else -> ask ta endpoint for list of devices & their properties & create a device info file
	// TODO: make sure to update the device file anytime somebody reserves or releases a device

	devicesList := handlers.DevicesMap{
		Devices: make(map[handlers.DeviceName]*modals.DeviceProps),
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		info, err := GetTmtWorkers(taProjectEndpoint)
		if err != nil {
			// No local device info and error while getting info from TMT -> fail program
			log.Fatal(err)
		}

		for _, worker := range info.Workers {
			devicesList.Devices[handlers.DeviceName(worker.Name)] = &modals.DeviceProps{
				Name:         worker.Name,
				Reserved:     false,
				ReservedBy:   "",
				ReservedTime: time.Now(), // TODO: how do i provide empty time??
				Props:        worker.Properties,
			}
		}

		devicesList.SynchronizeToFile() // save all obtained data to file
	} else if err != nil {
		// In case there is an error different that FileNotFound -> fail program
		log.Fatal(err)
	}

	// happy path
	log.Println(string(data))
	/*
		fileData, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Could not read users file %s", path)
		}

		err = json.Unmarshal(fileData, &usersList)
		if err != nil {
			log.Fatalf("Could not parse users file %s. Error: %+v", path, err)
		}
	*/

	/*
		devicesList = map[handlers.DeviceName]*modals.DeviceProps{
			"splinter":  &modals.DeviceProps{"splinter", false, "", time.Now()},
			"shredder":  &modals.DeviceProps{"shredder", false, "", time.Now()},
			"donatello": &modals.DeviceProps{"donatello", true, "Antanas", time.Now()},
		}
	*/

	log.Printf("device list loaded successfully\n%+v", devicesList)
	return devicesList
}
