package device

import (
	"testing"
)

func TestNewDevicesMapFromJson(t *testing.T) {
	json_input := []byte(`
            {
                "Devices":{
                    "donatello":{
                        "Name":"donatello",
                        "Reserved":false,
                        "AutoRelease":true,
                        "ReservedBy":"",
                        "ReservedById":"",
                        "ReservedTime":"0001-01-01T00:00:00Z",
                        "serial_proxy_host":"10.208.1.21",
                        "serial_proxy_port":5566,
                        "device_serial_number":"9QAA0776"
                     }
                 }
            }`)

	devicesMap := NewDevicesMapFromJson(json_input)

	if len(devicesMap.Devices) != 1 {
		t.Errorf("Expected 1 device to be process but got %d", len(devicesMap.Devices))
	}

	device, ok := devicesMap.Devices[DeviceName("donatello")]
	if !ok {
		t.Errorf("Could not find expected device 'donatello' in devicesMap =  %v", devicesMap.Devices)
	}

	if device.Reserved {
		t.Errorf("Expected 'donatello' to be free (it's not).")
	}

	if device.ReservedBy != "" || device.ReservedById != "" {
		t.Errorf(`Expected 'donatello' to not have any reserved info but it does:
		         ReservedBy: %v, ReservedById: %v.`, device.ReservedBy, device.ReservedById)
	}

	if !device.AutoRelease {
		t.Errorf("Expected 'donatello' to be marked for auto-release (it's not).")
	}

	if device.SerialProxyHost != "10.208.1.21" {
		t.Errorf("Expected 'donatello' proxy host to be '10.208.1.21' but got %v", device.SerialProxyHost)
	}

	if device.SerialProxyPort != 5566 {
		t.Errorf("Expected 'donatello' proxy port to be 5566 but got %v", device.SerialProxyPort)
	}

	if device.DeviceSerialNumber != "9QAA0776" {
		t.Errorf("Expected 'donatello' serial number to be '9QAA0776' but got %v", device.DeviceSerialNumber)
	}
}
