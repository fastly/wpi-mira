package common

type BGPMessage struct {
	Timestamp      float64
	BGPMessageType string
	PeerIP         string
	PeerASN        uint32
	Prefix         string
}

type Window struct {
	Filter    string
	BucketMap map[int64][]BGPMessage
}