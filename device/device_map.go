package device

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/AngelVI13/slack-assistant/config"
)

type DeviceName string

type DevicesMap struct {
	Devices map[DeviceName]*DeviceProps
	config  *config.Config
}

func NewDevicesMap() DevicesMap {
	return DevicesMap{
		Devices: make(map[DeviceName]*DeviceProps),
	}
}

// NewDevicesMapFromJson Takes json data as input and returns a populated DevicesMap object
func NewDevicesMapFromJson(data []byte, config *config.Config) DevicesMap {
	devicesList := NewDevicesMap()
	devicesList.synchronizeFromFile(data)
	devicesList.config = config
	return devicesList
}

func (d *DevicesMap) SynchronizeToFile() {
	data, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(d.config.DevicesFilename, data, 0666)
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

// GetDevicesInfo Returns DevicesInfo in a particular order. First are devices
// that have been reserved by provided user. Then are all the other devices
// sorted by status (reserved or not) and then sub-sorted by name.
func (d *DevicesMap) GetDevicesInfo(user string) DevicesInfo {
	// Group devices in 2 groups -> belonging to given user or not
	// The group that doesn't belong to user will be sorted by name and by status (reserved or not)
	userDevices := make(DevicesInfo, 0)
	nonUserDevices := make(DevicesInfo, 0)
	for _, d := range d.Devices {
		if d.Reserved && d.ReservedBy == user {
			userDevices = append(userDevices, d)
		} else {
			nonUserDevices = append(nonUserDevices, d)
		}
	}

	// NOTE: This sorts the device list starting from free devices
	sort.Slice(nonUserDevices, func(i, j int) bool {
		return !nonUserDevices[i].Reserved
	})

	firstTaken := -1 // Index of first taken device
	for i, device := range nonUserDevices {
		if device.Reserved {
			firstTaken = i
			break
		}
	}

	// NOTE: this might be unnecessary but it shows devices in predicable way in UI so its nice.
	// If all devices are free or all devices are taken, sort by name
	if firstTaken == -1 || firstTaken == 0 {
		sort.Slice(nonUserDevices, func(i, j int) bool {
			return nonUserDevices[i].Name < nonUserDevices[j].Name
		})
	} else {
		// split devices into 2 - free & taken
		// sort each sub slice based on device name/port
		free := nonUserDevices[:firstTaken]
		taken := nonUserDevices[firstTaken:]

		sort.Slice(free, func(i, j int) bool {
			return free[i].Name < free[j].Name
		})

		sort.Slice(taken, func(i, j int) bool {
			return taken[i].Name < taken[j].Name
		})
	}

	allDevices := make(DevicesInfo, 0, len(d.Devices))
	allDevices = append(allDevices, userDevices...)
	allDevices = append(allDevices, nonUserDevices...)
	return allDevices
}

func (d *DevicesMap) Reserve(deviceName, user, userId string, autoRelease bool) (err string) {
	device, ok := d.Devices[DeviceName(deviceName)]
	if !ok {
		log.Fatalf("Wrong device name %s, %+v", deviceName, d)
	}
	// Only inform user if it was someone else that tried to reserved his device.
	// This prevents an unnecessary message if you double clicked the reserve button yourself
	if device.Reserved && device.ReservedById != userId {
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
	device, ok := d.Devices[DeviceName(deviceName)]
	if !ok {
		log.Fatalf("Wrong device name %v, %+v", deviceName, d)
	}

	log.Printf("RELEASE: User (%s) released (%s) device.", user, deviceName)
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

	// Need to synchronize changes from file otherwise the state won't be preserved after restart
	// NOTE: This ends up synchronizing to file more than once since the function can be called
	// multiple times within the specified auto release hour (even if nothing has changed in the devices list).
	d.SynchronizeToFile()
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
	requestBodyJson, err := json.MarshalIndent(requestBody, "", "\t")

	if err != nil {
		log.Fatalf("error while marshalling request body [%v] - %v", requestBody, err)
	}

	// Specify timeout.
	// TODO: Investigate adding transport to client for setting up proxy
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	proxyUrl := fmt.Sprintf("%s/proxy", d.config.TaEndpoint)
	resp, err := client.Post(proxyUrl, "application/json", bytes.NewBuffer(requestBodyJson))

	if err != nil {
		return fmt.Sprintf("Restart proxy POST request to TA_ENDPOINT failed: err=%+v", err)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Could not read response body from restart proxy POST req: err=%+v", err)
	}

	cmdOutput := jsonPrettyPrint(responseBody)
	cmdOutput = fmt.Sprintf("```%s```", cmdOutput) // Display cmdOutput as code block for better readability

	return cmdOutput
}

func jsonPrettyPrint(in []byte) string {
	var out bytes.Buffer
	err := json.Indent(&out, in, "", "\t")
	if err != nil {
		return string(in)
	}
	return out.String()
}
