package device

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type DeviceName string
type AccessRight int

// NOTE: Currently access rights are not used
const (
	STANDARD AccessRight = iota
	ADMIN
)

type DevicesMap struct {
	Devices map[DeviceName]*DeviceProps
}

func NewDevicesMap() DevicesMap {
	return DevicesMap{
		Devices: make(map[DeviceName]*DeviceProps),
	}
}

// NewDevicesMapFromJson Takes json data as input and returns a populated DevicesMap object
func NewDevicesMapFromJson(data []byte) DevicesMap {
	devicesList := NewDevicesMap()
	devicesList.synchronizeFromFile(data)
	return devicesList
}

func (d *DevicesMap) SynchronizeToFile() {
	data, err := json.Marshal(d)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(os.Getenv("SL_DEVICES_FILE"), data, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("INFO: Wrote device list to file")
}

func (d *DevicesMap) synchronizeFromFile(data []byte) {
	// Unmarshal the provided data into the solid map
	err := json.Unmarshal(data, d)
	if err != nil {
		log.Fatalf("Could not parse devices file. Error: %+v", err)
	}
}

func (d *DevicesMap) Reserve(deviceName, user, userId string, autoRelease bool) (err string) {
	device, ok := d.Devices[DeviceName(deviceName)]
	if !ok {
		log.Fatalf("Wrong device name %s, %+v", deviceName, d)
	}
	if device.Reserved {
		reservedTime := device.ReservedTime.Format("Mon 15:04")
		return fmt.Sprintf("*Error*: Could not reserve *%s*. *%s* has just reserved it (at *%s*)", deviceName, device.ReservedBy, reservedTime)
	}
	log.Printf("RESERVE: User (%s) reserved device (%s) with auto release (%v)", user, deviceName, autoRelease)

	device.Reserved = true
	device.ReservedBy = user
	device.ReservedById = userId
	device.ReservedTime = time.Now()
	device.AutoRelease = autoRelease

	d.SynchronizeToFile()
	return ""
}

func (d *DevicesMap) Release(deviceName, user string) (victimId, err string) {
	log.Printf("RELEASE: User (%s) released (%s) device.", user, deviceName)

	device, ok := d.Devices[DeviceName(deviceName)]
	if !ok {
		log.Fatalf("Wrong device deviceName %s, %+v", deviceName, d)
	}

	device.Reserved = false
	d.SynchronizeToFile()

	if device.ReservedBy != user {
		return device.ReservedById, fmt.Sprintf(":warning: *%s* released your (*%s*) device (*%s*)", user, device.ReservedBy, device.Name)
	}
	return "", ""
}

func (d *DevicesMap) AutoRelease(when time.Time) {
	// Only release devices at the specified hour (hour is [0;23])
	now := time.Now()
	if now.Hour() != when.Hour() {
		return
	}

	for _, device := range d.Devices {
		if device.Reserved && device.AutoRelease {
			device.Reserved = false
			device.AutoRelease = false
		}
	}
}

func (d *DevicesMap) RestartProxies(deviceNames []string, user string) string {
	log.Printf("RESTART_PROXY: User (%s) restarted (%s) device/s.", user, deviceNames)

	// Check that device names selected by user are part of configured devices
	// NOTE: this should technically never fails because the user is only showed
	// 	 valid device names but maybe keep this check here just in case ???
	for _, deviceName := range deviceNames {
		_, ok := d.Devices[DeviceName(deviceName)]
		if !ok {
			log.Fatalf("Wrong device deviceName %s, %+v", deviceName, d)
		}
	}

	requestBody := map[string]string{
		"command":      "restart",
		"device_names": strings.Join(deviceNames, ","),
	}
	requestBodyJson, err := json.Marshal(requestBody)

	if err != nil {
		log.Fatalf("error while marshalling request body [%v] - %v", requestBody, err)
	}

	proxyUrl := fmt.Sprintf("%s/proxy", os.Getenv("SL_TA_ENDPOINT"))
	resp, err := http.Post(proxyUrl, "application/json", bytes.NewBuffer(requestBodyJson))

	if err != nil {
		log.Fatal(err)
	}

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)

	fmt.Println(res)
	// TODO: do something with the response
	// TODO: might have to do this POST request asyncronously cause slack is expecting
	// configurmation at some point. Maybe send the response back to the user as a DM?
	return ""
}
