package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"
	"strings"
	"strconv"
	"log"
	"path/filepath"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/plotter"
	"math"
	"sort"
	//"image/color"
)

/* Examples:
BGP4MP|1467212398|A|206.223.119.67|852|179.124.136.0/21|852 6453 6762 52965|IGP|206.223.119.67|0|3851||NAG||
BGP4MP|1467180898|A|2001:504:0:4::852:1|852|2001:df0:230::/48|852 174 6762 7713 56237|IGP|2001:504:0:4::852:1|0|4011||NAG||
BGP4MP|1467180898|A|206.223.119.67|852|182.23.47.0/24|852 6453 4800|IGP|206.223.119.67|0|4011||NAG||
BGP4MP|1467180898|A|206.223.119.67|852|182.23.47.0/24|852 174 2914 4800|IGP|206.223.119.67|0|0||NAG||
BGP4MP|1467180900|A|206.223.119.67|852|179.124.136.0/21|852 174 6762 52965|IGP|206.223.119.67|0|3851||NAG||
BGP4MP|1467180900|A|206.223.119.67|852|168.235.164.0/24|852 7922 33662 63483 63483 63483 63483|IGP|206.223.119.67|0|0||NAG|63483 50.233.78.210|
BGP4MP|1467180900|A|206.223.119.67|852|168.235.174.0/24|852 7922 33662 63483 63483 63483 63483|IGP|206.223.119.67|0|0||NAG|63483 50.233.78.210|
BGP4MP|1467180900|A|206.223.119.67|852|168.235.173.0/24|852 7922 33662 63483 63483 63483 63483|IGP|206.223.119.67|0|0||NAG|63483 50.233.78.210|
BGP4MP|1467180900|A|206.223.119.67|852|168.235.172.0/24|852 7922 33662 63483 63483 63483 63483 4|IGP|206.223.119.67|0|0||NAG|63483 50.233.78.210|
BGP4MP|1467180900|A|206.223.119.67|852|168.235.171.0/24|852 7922 33662 63483 63483 63483|IGP|206.223.119.67|0|0||NAG|63483 50.233.78.210|
BGP4MP|1467180900|A|206.223.119.67|852|66.194.175.0/24|852 7922 33662 63483 63483 63483|IGP|206.223.119.67|0|0||NAG|63483 50.233.78.210|
BGP4MP|1467180900|A|206.223.119.67|852|66.194.174.0/24|852 7922 33662 63483 63483 63483|IGP|206.223.119.67|0|0||NAG|63483 50.233.78.210|
*/

/*
type BGPMessage struct {
	MessageType      string // BGP4MP
	Timestamp        int64 // 1467180900
	BGPMessageType   string // A
	PeerIP           string // 206.223.119.67
	PeerASN          uint32 // 852
	Prefix           string // 66.194.174.0/24
	ASPath           []uint32 // 852 7922 33662 63483 63483 63483
	Origin           string // IGP
	NextHopIP        string // 206.223.119.67
	LocalPref        uint32 // 0
	MED              uint32 // 0
	Community        string // 
	AtomicAggregator string // NAG
	Aggregator       string // 63483 50.233.78.210
}
*/

type BGPMessage struct {
	Timestamp        float64 // 1467180900
	BGPMessageType   string // A
	PeerIP           string // 206.223.119.67
	PeerASN          uint32 // 852
	Prefix           string // 66.194.174.0/24
}

func main() {
	startTime := time.Now()

	folderDir := "bgpfilesALOT"

	// Gets all files from directory we want to test
	files, err := os.ReadDir("testdata/"+folderDir)
	if err != nil {
		log.Fatal(err)
	}

	// Array for storing all BGPMessages
	var allBGPMessages []BGPMessage

	// Iterate through the files in the directory
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".bz2" {
			filePath := filepath.Join("testdata/" + folderDir, file.Name())
			bgpMessages, err := parseBGPFile(filePath)
			if err != nil {
				fmt.Println("Error parsing file:", err)
				continue
			}
			allBGPMessages = append(allBGPMessages, bgpMessages...)
		}
	}

	// Creates map of total frequencies for timestamp
	freqMap := buildFrequencyMap(allBGPMessages)
	
	//fmt.Println(freqMap)

	// Smallest timestamp is the first BGPMessage we store, largest is the last one
	smallestTimestamp := int64(allBGPMessages[0].Timestamp)
	largestTimestamp := int64((allBGPMessages[len(allBGPMessages) - 1]).Timestamp)

	// Round timestamps down to nearest multiple of 60 because decimals are annoying 
	smallestTimestamp = int64(math.Floor(float64(smallestTimestamp)/60.0) * 60)
	largestTimestamp = int64(math.Floor(float64(largestTimestamp)/60.0) * 60)

	// Going to store frequencies in one minute time buckets instead of per second
    bucketMap := make(map[int64]int)

    currentTimestamp := smallestTimestamp

    // Initializes minute long bins in map
    for currentTimestamp <= largestTimestamp {
        bucketMap[currentTimestamp] = 0
        currentTimestamp += 60 // Increment by 60 seconds
    }

    for timestamp, frequency := range freqMap {
        // Get corresponding bucket and then add frequency to it
        bucketTimestamp := (timestamp / 60) * 60
        bucketMap[bucketTimestamp] += frequency
    }

	bucketTotal := 0
    for _, frequency := range bucketMap {
		bucketTotal += frequency
    }

	// Gets sorted array of timestamps to feed into Analysis
	sortedFrequencies := getSortedFrequencies(bucketMap)

	fmt.Println("Sorted Frequencies: ", sortedFrequencies)
	fmt.Println("BucketMap: ", bucketMap)
	fmt.Println("Sum of frequencies: ", bucketTotal)
	fmt.Println("Number of buckets: ", len(bucketMap))
	fmt.Printf("Timestamp Difference: %d\n", (largestTimestamp - smallestTimestamp) / 60)

	plotData(bucketMap, folderDir)

	elapsedTime := time.Since(startTime)
	fmt.Println("Total execution time: ", elapsedTime)
	fmt.Println("Message count: ", len(allBGPMessages))

}

func parseBGPFile(filePath string) ([]BGPMessage, error) {
    cmd := exec.Command("bgpdump", "-m", "-O", "tempbgpdump.txt", filePath)
    err := cmd.Run()
    if err != nil {
        return nil, err
    }

    file, err := os.Open("tempbgpdump.txt")
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var bgpMessages []BGPMessage
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := scanner.Text()
        //fmt.Println(line)
        bgpMsg, err := parseBGPMessage(line)
        if err != nil {
            fmt.Printf("Error parsing message: %v\n", err)
            continue // Skip the invalid message
        }
        bgpMessages = append(bgpMessages, bgpMsg)
    }

	os.Remove("tempbgpdump.txt")

    if scanner.Err() != nil {
        return nil, scanner.Err()
    }

    return bgpMessages, nil
}


func plotData(bucketMap map[int64]int, folderDir string) {
	p:= plot.New()

    p.Title.Text = "Frequency vs. Timestamp"
    p.X.Label.Text = "Timestamp"
    p.Y.Label.Text = "Frequency"

    scatterData := make(plotter.XYs, 0)

    for timestamp, freq := range bucketMap {
        scatterData = append(scatterData, struct{ X, Y float64 }{float64(timestamp), float64(freq)})
    }

    scatter, err := plotter.NewScatter(scatterData)
    if err != nil {
        panic(err)
    }

    p.Add(scatter)

    if err := p.Save(8*vg.Inch, 8*vg.Inch, "plots/" + folderDir + ".png"); err != nil {
        panic(err)
    }
}

func buildFrequencyMap(BGPMessages []BGPMessage) map[int64]int {
	freqMap := make(map[int64]int)
	for _, msg := range BGPMessages {
		freqMap[int64(msg.Timestamp)] += 1
	}

	return freqMap
}


/*
func plotData(bucketMap map[int64]int, folderDir string) {
	p := plot.New()

	p.Title.Text = "Frequency vs. Timestamp"
	p.X.Label.Text = "Timestamp"
	p.Y.Label.Text = "Frequency"

	// Create a scatter plot with red points
	pts := make(plotter.XYs, 0)
	for timestamp, freq := range bucketMap {
		pts = append(pts, struct{ X, Y float64 }{float64(timestamp), float64(freq)})
	}
	s, err := plotter.NewScatter(pts)
	if err != nil {
		panic(err)
	}

	s.GlyphStyle.Color = color.RGBA{255, 0, 0, 255} // Red color
	//s.GlyphStyle.Color = color.RGBA{255, 0, 0, 0} // Light blue color
	p.Add(s)

	// Save the plot
	if err := p.Save(8*vg.Inch, 8*vg.Inch, "plots/"+folderDir+".png"); err != nil {
		panic(err)
	}
}
*/

func getSortedFrequencies(bucketMap map[int64]int) []int {
	var timestamps []int64

	// Create a slice of timestamps and a corresponding slice of values in the order of timestamps
	for timestamp, _ := range bucketMap {
		timestamps = append(timestamps, timestamp)
	}

	// Sort the timestamps in ascending order
	sort.Slice(timestamps, func(i, j int) bool {
		return timestamps[i] < timestamps[j]
	})

	// Create a new slice of values in the order of sorted timestamps
	sortedValues := make([]int, len(timestamps))
	for i, timestamp := range timestamps {
		sortedValues[i] = bucketMap[timestamp]
	}
	
	return sortedValues
}

func parseBGPMessage(data string) (BGPMessage, error) {
    fields := strings.Split(data, "|")

    // If message type is withdrawal (W) -> there's only the first 6 fields
    if len(fields) == 6 && fields[2] == "W" {
		return BGPMessage {
			Timestamp:      parseTimestamp(fields[1]),
            BGPMessageType: fields[2],
            PeerIP:         fields[3],
            PeerASN:        parseUint32(fields[4]),
            Prefix:         fields[5],
		}, nil
		/*
        return BGPMessage {
            MessageType:    fields[0],
            Timestamp:      parseTimestamp(fields[1]),
            BGPMessageType: fields[2],
            PeerIP:         fields[3],
            PeerASN:        parseUint32(fields[4]),
            Prefix:         fields[5],
        }, nil
		*/
    } else if len(fields) == 15 && fields[2] == "A" {
        // If the message is advertisement (A) -> parse like normal
		return BGPMessage {
			Timestamp:      parseTimestamp(fields[1]),
            BGPMessageType: fields[2],
            PeerIP:         fields[3],
            PeerASN:        parseUint32(fields[4]),
            Prefix:         fields[5],
		}, nil
		/*
        return BGPMessage{
            MessageType:      fields[0],
            Timestamp:        parseTimestamp(fields[1]),
            BGPMessageType:   fields[2],
            PeerIP:           fields[3],
            PeerASN:          parseUint32(fields[4]),
            Prefix:           fields[5],
            ASPath:           parseASPath(fields[6]),
            Origin:           fields[7],
            NextHopIP:        fields[8],
            LocalPref:        parseUint32(fields[9]),
            MED:              parseUint32(fields[10]),
            Community:        fields[11],
            AtomicAggregator: fields[12],
            Aggregator:       fields[13],
        }, nil
		*/
    }

    return BGPMessage{}, fmt.Errorf("Invalid BGP message: %s", data)
}

func parseTimestamp(timestampStr string) float64 {
	timestamp, err := strconv.ParseFloat(timestampStr, 64)
	if err != nil {
		// Handle the error if the conversion fails
		fmt.Println("Error Parsing Float from Timestamp: ", err)
		return 0.0 // or any other appropriate value
	}
	return timestamp
}

func parseUint32(valueStr string) uint32 {
	value, _ := strconv.ParseUint(valueStr, 10, 32)
	return uint32(value)
}

func parseASPath(asPathStr string) []uint32 {
	parts := strings.Fields(asPathStr)
	asPath := make([]uint32, len(parts))
	for i, part := range parts {
		asPath[i] = parseUint32(part)
	}
	return asPath
}
