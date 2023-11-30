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

type Window struct {
	Filter    string
	BucketMap map[time.Time][]BGPMessage
}
