package main

type resolution uint8

const (
	resolutionDaily resolution = iota
	resolutionHourly
	resolution10Minutes
)

func (r resolution) String() string {
	switch r {
	case resolutionDaily:
		return "daily"
	case resolutionHourly:
		return "hourly"
	case resolution10Minutes:
		return "10_minutes"
	default:
		return "unsupported"
	}
}
