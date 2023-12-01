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

type Result struct {
	AllOutliers []OutlierInfo         `json:"AllOutliers,omitempty"` //list of all the outliers
	AllFreq     map[time.Time]float64 `json:"AllFreq,omitempty"`     //map of all the frequencies by timestamp to avoid repeats and keep track of missed data
}

type OutlierInfo struct {
	Timestamp time.Time `json:"Timestamp"`           //time stamps of the outlier bucket
	Algorithm int       `json:"Algorithm,omitempty"` //0 if mad 1 if shakeAlert 2 if both
	Count     float64   `json:"Count,omitempty" `    //the number of messages in the bucket that is an outlier
	//MsgsFile  *os.File  `json:"MsgsFile,omitempty"`  //`json:"MsgsFile,omitempty"`  a file that contains all the messages from the outlier buckets
}
