package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/AngelVI13/slack-assistant/device"
)

type Worker struct {
	Name        string             `json:"name"`
	Password    string             `json:"password"`
	WorkerSetup []string           `json:"worker_setup"`
	Properties  device.WorkerProps `json:"properties"`
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

	log.Printf("INIT: Workers info fetched successfully")
	return &workers, nil
}

func GetDevices(path, taProjectEndpoint string) device.DevicesMap {
	// 1.a If device info file exists -> read info from there
	// 1.b Else -> ask TA endpoint for list of devices & their properties & create a device info file
	var devicesList device.DevicesMap

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		info, err := GetTmtWorkers(taProjectEndpoint)
		if err != nil {
			// No local device info and error while getting info from TMT -> fail program
			log.Fatal(err)
		}

		devicesList = device.NewDevicesMap()

		for _, worker := range info.Workers {
			devicesList.Devices[device.DeviceName(worker.Name)] = &device.DeviceProps{
				Name: worker.Name,
				ReservedProps: device.ReservedProps{
					Reserved: false,
				},
				WorkerProps: worker.Properties,
			}
		}

		devicesList.SynchronizeToFile() // save all obtained data to file
	} else if err != nil {
		// In case there is an error different that FileNotFound -> fail program
		log.Fatal(err)
	} else {
		// Devices file exists -> read from there
		devicesList = device.NewDevicesMapFromJson(data)
	}

	loadedDeviceNum := len(devicesList.Devices)
	if loadedDeviceNum == 0 {
		log.Fatalf("No devices found in (%s).", path)
	}

	log.Printf("INIT: Device list loaded successfully (%d devices configured)", loadedDeviceNum)
	return devicesList
}
