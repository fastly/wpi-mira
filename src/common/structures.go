package common

import (
	"time"
)

type BGPMessage struct {
	Timestamp      time.Time
	BGPMessageType string
	PeerIP         string
	PeerASN        uint32
	Prefix         string
}

type Window struct {
	Filter    string
	BucketMap map[time.Time][]BGPMessage
}
