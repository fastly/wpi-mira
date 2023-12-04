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

type Result struct {
	Filter      string                `json:"Filter"`                //filter
	AllOutliers []OutlierInfo         `json:"AllOutliers,omitempty"` //list of all the outliers
	AllFreq     map[time.Time]float64 `json:"AllFreq,omitempty"`     //map of all the frequencies by timestamp to avoid repeats and keep track of missed data
}

type OutlierInfo struct {
	Timestamp time.Time `json:"Timestamp"`           //time stamps of the outlier bucket
	Algorithm int       `json:"Algorithm,omitempty"` //0 if mad 1 if shakeAlert 2 if both
	Count     float64   `json:"Count,omitempty" `    //the number of messages in the bucket that is an outlier
	//MsgsFile  *os.File  `json:"MsgsFile,omitempty"`  //`json:"MsgsFile,omitempty"`  a file that contains all the messages from the outlier buckets
}

//result struct needed to simplify getting all the info needed to be dispayed
//use json marshall to parse and simplify the code to store
/*type Result struct {
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
*/
