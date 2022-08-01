package parking

import (
	"fmt"

	"github.com/AngelVI13/slack-assistant/props"
)

type ParkingSpace struct {
	Number int
	Floor  int
	props.ReservedProps
}

type ParkingSpaces map[int]*ParkingSpace

func (p *ParkingSpace) GetPropsText() string {
	return fmt.Sprintf("(%d floor)", p.Floor)
}

func (p *ParkingSpace) GetStatusEmoji() string {
	emoji := ":large_green_circle:"
	if p.Reserved {
		emoji = ":large_orange_circle:"
	}
	return emoji
}

// GetStatusDescription Get device status description i.e. reserved, by who, when, etc.
// Returns empty string if device is free
func (p *ParkingSpace) GetStatusDescription() string {
	status := ""
	if p.Reserved {
		timeStr := p.ReservedTime.Format("Mon 15:04")
		status = fmt.Sprintf("_:bust_in_silhouette:*%s*\ton\t:clock1: *%s*_", p.ReservedBy, timeStr)
	}
	return status
}
