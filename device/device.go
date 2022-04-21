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

func (p *DeviceProps) GetStatusEmoji() string {
	emoji := ":large_green_circle:"
	if p.Reserved {
		emoji = ":large_orange_circle:"
	}
	return emoji
}

// GetStatusDescription Get device status description i.e. reserved, by who, when, etc.
// Returns empty string if device is free
func (p *DeviceProps) GetStatusDescription() string {
	status := ""
	if p.Reserved {
		timeStr := p.ReservedTime.Format("Mon 15:04")
		autoStatus := ""
		if p.AutoRelease {
			autoStatus = "\t:eject: *Auto*"
		}
		status = fmt.Sprintf("_:bust_in_silhouette:*%s*\ton\t:clock1: *%s*%s_", p.ReservedBy, timeStr, autoStatus)
	}
	return status
}
