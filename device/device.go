package device

import (
	"fmt"
	"time"
)

type WorkerProps struct {
	SerialProxyHost    string `json:"serial_proxy_host"`
	SerialProxyPort    int    `json:"serial_proxy_port"`
	DeviceSerialNumber string `json:"device_serial_number"`
}

type ReservedProps struct {
	Reserved     bool
	AutoRelease  bool
	ReservedBy   string
	ReservedById string
	ReservedTime time.Time
}

type DeviceProps struct {
	Name string
	ReservedProps
	WorkerProps
}

type DevicesInfo []*DeviceProps

func (p *DeviceProps) GetPropsText() string {
	return fmt.Sprintf(
		"Port: %d\tHost: %s\tS/N: %s",
		p.SerialProxyPort,
		p.SerialProxyHost,
		p.DeviceSerialNumber,
	)
}
