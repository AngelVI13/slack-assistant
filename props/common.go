package props

import "time"

type ReservedProps struct {
	Reserved     bool
	AutoRelease  bool
	ReservedBy   string
	ReservedById string
	ReservedTime time.Time
}
