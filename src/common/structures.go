package common

import (
	"net/netip"
	"time"
)

const (
	AnnouncementType = "A"
	WithdrawalType   = "W"
)

type BGPMessage struct {
	Timestamp      time.Time
	BGPMessageType string
	PeerIP         netip.Addr
	PeerASN        uint32
	Prefix         netip.Prefix
}

type Window struct {
	Filter    string
	BucketMap map[time.Time][]BGPMessage
}

//result struct needed to simplify getting all the info needed to be dispayed
//use json marshall to parse and simplify the code to store
type Result struct {
	Frequencies          []float64   `json:"Frequencies"`
	MADOutliers          []float64   `json:"MADOutliers"`
	MADTimestamps        []time.Time `json:"MADTimestamps"`
	ShakeAlertOutliers   []float64   `json:"ShakeAlertOutliers"`
	ShakeAlertTimestamps []time.Time `json:"ShakeAlertTimestamps"`
}
