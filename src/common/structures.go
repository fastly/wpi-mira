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

type Message struct {
	BGPMessage BGPMessage
	Filter     string
}

//used to write the struct onto a json output
type OutlierMessages struct {
	MADOutlierMessages [][]BGPMessage `json:"MADOutlierMessages"`
	ShakeAlertMessages [][]BGPMessage `json:"ShakeAlertMessages"`
}

type Window struct {
	Filter    string
	BucketMap map[time.Time][]BGPMessage
}

//result struct needed to simplify getting all the info needed to be dispayed
//use json marshall to parse and simplify the code to store
type Result struct {
	//config parameters
	Prefix     string `json:"Prefix"`
	ASN        string `json:"ASN"`
	PeerIP     string `json:"PeerIP"`
	WindowSize string `json:"WindowSize"`

	//outputs from processing the live messages
	Frequencies []float64 `json:"Frequencies"`

	MADOutliers   []float64   `json:"MADOutliers"`
	MADTimestamps []time.Time `json:"MADTimestamps"`

	ShakeAlertOutliers   []float64   `json:"ShakeAlertOutliers"`
	ShakeAlertTimestamps []time.Time `json:"ShakeAlertTimestamps"`
}
