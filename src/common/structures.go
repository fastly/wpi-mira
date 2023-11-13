package common

import (
	"net/netip"
	"time"
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
