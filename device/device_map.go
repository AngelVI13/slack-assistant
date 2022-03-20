package device

import (
    "encoding/json"
    "time"
    "log"
    "os"
    "fmt"
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

func (d *DevicesMap) Reserve(deviceName, user, userId string) (err string) {
	device, ok := d.Devices[DeviceName(deviceName)]
	if !ok {
		log.Fatalf("Wrong device name %s, %+v", deviceName, d)
	}
	if device.Reserved {
		reservedTime := device.ReservedTime.Format("Mon 15:04")
		return fmt.Sprintf("*Error*: Could not reserve *%s*. *%s* has just reserved it (at *%s*)", deviceName, device.ReservedBy, reservedTime)
	}
	log.Printf("RESERVE: User (%s) reserved device (%s)", user, deviceName)

	device.Reserved = true
	device.ReservedBy = user
	device.ReservedById = userId
	device.ReservedTime = time.Now()

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

