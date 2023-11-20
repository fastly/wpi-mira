package main

import (
	"BGPAlert/analyze"
	"BGPAlert/blt_mad"
	"BGPAlert/common"
	"BGPAlert/config"
	"BGPAlert/optimization"
	"bufio"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/netip"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Folder struct {
	FolderName string
	Urls       []string
}

func main() {
	// Specify the target website URL
	//to get the bgp updates the urls are of this format; sorted by the year and month
	//"http://routeviews.org/route-views.ny/bgpdata/2021.11/UPDATES/"
	//the links to .bz2 files contained within the main link in the configuration file; if the links can not be attained from the main link no files will be downloaded
	configStruct, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}
	config.ValidDateConfiguration(configStruct)
	mainUrl := configStruct.URLStaticData
	fmt.Println(mainUrl)

	// Make an HTTP GET request to the website
	resp, err := http.Get(mainUrl)
	if err != nil {
		fmt.Printf("Error making HTTP request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Parse the HTML content of the response
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Printf("Error parsing HTML: %v\n", err)
		os.Exit(1)
	}

	// Extract URLs from the HTML document
	partialUrls := extractBZ2URLs(doc)
	allFolders := createFolders(partialUrls, mainUrl, 10)
	for _, folder := range allFolders {
		err := downloadFolder(folder)
		if err != nil {
			return
		}
	}

	for i := 0; i <= len(allFolders)-1; i++ {
		outFile := fmt.Sprintf("static_data/rawBGPData/bgpTest%d.txt", i)
		outMinReqFile := fmt.Sprintf("static_data/minReq/bgpMinOutliers97thPercentileTest%d.txt", i)
		runProcessThroughOneBGPFolder(i, outFile, outMinReqFile)
	}
	//get all the means,mads, and taus into text files
	GetMadsMediansTausIntoTxt()

}

func downloadFile(url, folderPath string) error {
	// Create the folder if it doesn't exist
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return fmt.Errorf("error creating folder: %v", err)
	}

	// Extract the file name from the URL
	fileName := filepath.Base(url)

	// Create the file in the specified folder
	filePath := filepath.Join(folderPath, fileName)
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer out.Close()

	// Make the HTTP request to download the file
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %v", err)
	}
	defer response.Body.Close()

	// Check if the request was successful (status code 200)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error downloading file, status code: %d", response.StatusCode)
	}

	// Copy the contents of the response body to the file
	_, err = io.Copy(out, response.Body)
	if err != nil {
		return fmt.Errorf("error copying file contents: %v", err)
	}

	fmt.Printf("File downloaded successfully to: %s\n", filePath)
	return nil
}

func downloadFolder(folder Folder) error {
	for _, url := range folder.Urls {
		if err := downloadFile(url, folder.FolderName); err != nil {
			fmt.Printf("Error downloading file from URL %s: %v\n", url, err)
			// Continue with the next URL even if one fails
		}
	}
	return nil
}

//this works
func createFullUrls(urls []string, mainURL string) []string {
	var result []string

	for _, element := range urls {
		result = append(result, mainURL+element)
	}

	return result
}

func getIntervalSlices(list []string, interval int) [][]string {
	var result [][]string

	for i := 0; i < len(list); i += interval {
		end := i + interval
		if end > len(list) {
			end = len(list)
		}
		result = append(result, list[i:end])
	}

	return result
}

func createFolders(urls []string, mainURL string, numFiles int) []Folder {
	var folderList []Folder
	fullUrlsList := createFullUrls(urls, mainURL)
	setsOfUrls := getIntervalSlices(fullUrlsList, numFiles)

	for i := 0; i < len(setsOfUrls); i++ {
		folderName := fmt.Sprintf("static_data/bgpTest%d", i)

		// Create a Folder struct and append it to the list
		folder := Folder{
			FolderName: folderName,
			Urls:       setsOfUrls[i],
		}
		folderList = append(folderList, folder)
	}
	return folderList
}

func extractBZ2URLs(n *html.Node) []string {
	var urls []string

	if n.Type == html.ElementNode && (n.Data == "a" || n.Data == "img" || n.Data == "link") {
		for _, attr := range n.Attr {
			if attr.Key == "href" || attr.Key == "src" {
				if strings.HasSuffix(attr.Val, ".bz2") {
					urls = append(urls, attr.Val)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		urls = append(urls, extractBZ2URLs(c)...)
	}

	return urls
}

func GetMadsMediansTausIntoTxt() {
	folderPath := "static_data/rawBGPData"
	allMads := []float64{}
	allMedians := []float64{}
	allTaus := []float64{}

	// Read all files in the folder
	files, _ := ioutil.ReadDir(folderPath)

	// Count the number of files (excluding directories)
	fileCount := 0
	for _, file := range files {
		if !file.IsDir() {
			fileCount++
		}
	}

	//iterate through all the files and append the mads into one array
	for i := 0; i < fileCount; i++ {
		inputFilePath := fmt.Sprintf("static_data/rawBGPData/bgpTest%d.txt", i)
		allMads = append(allMads, findMadOfTextFile(inputFilePath))
	}
	//iterate through all the files and append the medians into one array
	for i := 0; i < fileCount; i++ {
		inputFilePath := fmt.Sprintf("static_data/rawBGPData/bgpTest%d.txt", i)
		allMedians = append(allMedians, findMedianOfTextFile(inputFilePath))
	}

	//iterate through all the files and append the optimal taus into one array
	for i := 0; i < fileCount; i++ {
		inputFilePath := fmt.Sprintf("static_data/rawBGPData/bgpTest%d.txt", i)
		inputMinFilePath := fmt.Sprintf("static_data/minReq/bgpMinOutliers97thPercentileTest%d.txt", i)
		allTaus = append(allTaus, findTauOfTextFile(inputFilePath, inputMinFilePath))
	}

	//write the array into a text file
	blt_mad.SaveArrayToFile("static_data/madsFound.txt", allMads)
	blt_mad.SaveArrayToFile("static_data/mediansFound.txt", allMedians)
	blt_mad.SaveArrayToFile("static_data/tausFound.txt", allTaus)
}

func findTauOfTextFile(inputTextFile string, minReqTestFile string) float64 {
	arrMain, _ := blt_mad.TxtIntoArrayFloat64(inputTextFile)
	arrMinReq, _ := blt_mad.TxtIntoArrayFloat64(minReqTestFile)
	//fmt.Println(arrMinReq)
	tau := optimization.FindTauForMinReqOutput(arrMain, arrMinReq)
	return tau
}

func findMedianOfTextFile(inputTextFile string) float64 {
	arr, _ := blt_mad.TxtIntoArrayFloat64(inputTextFile)
	median := blt_mad.FindMedian(arr)
	return median
}

func findMadOfTextFile(inputTextFile string) float64 {
	arr, _ := blt_mad.TxtIntoArrayFloat64(inputTextFile)
	mean := blt_mad.Mad(arr)
	return mean
}

//modify the function to write outputs onto a file
func AnalyzeBGPMessagesWriteOntoFile(windowChannel chan []common.Window, freqOutFile string, minReqFile string) {
	for windows := range windowChannel {
		for _, w := range windows {
			//fmt.Println("Received Window: ")
			bucketMap := w.BucketMap

			// Convert BucketMap to a map of timestamp to length of messages
			lengthMap := make(map[time.Time]float64)

			for timestamp, messages := range bucketMap {
				lengthMap[timestamp] = float64(len(messages))
			}

			// Turn map into sorted array of frequencies by timestamp
			sortedFrequenciesFloatArray := analyze.GetSortedFrequencies(lengthMap)

			// Convert to float array for analysis functions
			//floatArray := int64ArrayToFloat64Array(sortedFrequencies)
			//save frequency array into the freqOutFile
			blt_mad.SaveArrayToFile(freqOutFile, sortedFrequenciesFloatArray)

			//get min reqArray for the 97th percentile
			minReqArray := blt_mad.GetValuesLargerThanPercentile(sortedFrequenciesFloatArray, 97)
			blt_mad.SaveArrayToFile(minReqFile, minReqArray)
		}
	}
}

func runProcessThroughOneBGPFolder(num int, outFile string, outMinReqFile string) {
	// WaitGroup for waiting on goroutines to finish
	var wg sync.WaitGroup

	// Channel for sending BGP messages between parsing and processing
	msgChannel := make(chan []common.BGPMessage)

	// Channel for sending windows from processing to analyzing
	windowChannel := make(chan []common.Window)

	// Start the goroutines
	wg.Add(3)
	// Can change folder directory to any folder inside of src/staticdata
	inputFolderPath := fmt.Sprintf("bgpTest%d", num)
	//outFile := fmt.Sprintf("bgpTest%d.txt", num)

	go func() {
		parseStaticFile(inputFolderPath, msgChannel)
		wg.Done()
	}()

	go func() {
		processBGPMessagesStatic(msgChannel, windowChannel)
		wg.Done()
	}()

	go func() {
		AnalyzeBGPMessagesWriteOntoFile(windowChannel, outFile, outMinReqFile)
		wg.Done()
	}()

	wg.Wait()

}

// Reads in static files from a directory, parses them into BGPMessage struct and puts them into channel for processor to process
func parseStaticFile(folderDir string, msgChannel chan []common.BGPMessage) {
	const directoryPath = "static_data/"

	// Gets all files from directory we want to test
	files, err := os.ReadDir(directoryPath + folderDir)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Iterate through the files in the directory
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".bz2" {
			filePath := filepath.Join(directoryPath, folderDir, file.Name())
			bgpMessages, err := parseBGPFile(filePath)
			if err != nil {
				fmt.Println("Error parsing file:", err)
				continue
			}

			// Send the parsed BGP messages to the channel
			msgChannel <- bgpMessages
		}
	}

	// Close the channel to signal that all messages have been sent
	close(msgChannel)
}

// Runs BGPDump on individual static bgp files and extracts all BGPMessages from that individual file
func parseBGPFile(filePath string) ([]common.BGPMessage, error) {
	const tempFileName = "tempbgpdump.txt"

	cmd := exec.Command("bgpdump", "-m", "-O", tempFileName, filePath)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(tempFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var bgpMessages []common.BGPMessage
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		bgpMsg, err := parseBGPMessage(line)
		if err != nil {
			fmt.Printf("Error parsing message: %v\n", err)
			continue // Skip the invalid message
		}
		bgpMessages = append(bgpMessages, bgpMsg)
	}

	os.Remove(tempFileName)

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return bgpMessages, nil
}

// Takes a line of BGPDumped data and turns it into the BGPMessage struct
func parseBGPMessage(data string) (common.BGPMessage, error) {
	const announcementFields = 15
	const announcmentType = "A"
	const withdrawalFields = 6
	const withdrawalType = "W"

	fields := strings.Split(data, "|")

	timestamp, err := parseTimestamp(fields[1])
	if err != nil {
		return common.BGPMessage{}, fmt.Errorf("error parsing timestamp: %v", err)
	}

	bgpMessageType := fields[2]

	peerIP, err := netip.ParseAddr(fields[3])
	if err != nil {
		return common.BGPMessage{}, fmt.Errorf("error parsing address: %v", err)
	}

	peerASN, err := parseUint32(fields[4])
	if err != nil {
		return common.BGPMessage{}, fmt.Errorf("error parsing peer asn: %v", err)
	}

	prefix, err := netip.ParsePrefix(fields[5])
	if err != nil {
		return common.BGPMessage{}, fmt.Errorf("error parsing prefix: %v", err)
	}

	// Announcment examples: (Has 15 fields and messageType must be "A")
	// BGP4MP_ET|1638317699.191812|A|206.82.104.185|398465|102.66.116.0/24|398465 5713 37457 37457 37457 37457 37457 37457 37457 328471 328471 328471|IGP|206.82.104.185|0|0|5713:800 65101:1085 65102:1000 65103:276 65104:150|NAG|4200000002 10.102.100.2|
	// BGP4MP_ET|1638317699.744853|A|2001:504:36::6:1481:0:1|398465|2a10:cc42:1bb9::/48|398465 174 1299 20473|IGP|2001:504:36::6:1481:0:1|0|0|174:21000 174:22003|NAG||
	// BGP4MP_ET|1638317700.043880|A|206.82.104.185|398465|102.66.116.0/24|398465 30844 328471 328471 328471|IGP|206.82.104.185|0|0|30844:27 65101:1082 65102:1000 65103:276 65104:150|NAG|4200000002 10.102.100.2|

	// Withdrawal examples: (Has 6 fields and third field messageType must be "W")
	// BGP4MP_ET|1638317694.706880|W|2001:504:36::6:1481:0:1|398465|2a10:cc42:131c::/48
	// BGP4MP_ET|1638317679.511516|W|2001:504:36::6:1481:0:1|398465|2804:2688::/33

	if (len(fields) == announcementFields && bgpMessageType == announcmentType) || (len(fields) == withdrawalFields && bgpMessageType == withdrawalType) {
		return common.BGPMessage{
			Timestamp:      timestamp,
			BGPMessageType: bgpMessageType,
			PeerIP:         peerIP,
			PeerASN:        peerASN,
			Prefix:         prefix,
		}, nil
	}

	return common.BGPMessage{}, fmt.Errorf("invalid bgp message: %s", data)
}

// Converts a string timestamp into a time.Time object
func parseTimestamp(timestampStr string) (time.Time, error) {
	timestamp, err := strconv.ParseFloat(timestampStr, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing float from timestamp: %v", err)
	}

	seconds := int64(timestamp)
	nanoseconds := int64((timestamp - float64(seconds)) * 1e9)

	time := time.Unix(seconds, nanoseconds)

	return time, nil
}

// Converts a string uint and converts it into an uint32
func parseUint32(valueStr string) (uint32, error) {
	value, err := strconv.ParseUint(valueStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("error parsing uint32: %v", err)
	}

	return uint32(value), nil
}

// Constantly reads channel of messages and stores them in Window objects to send through windowChannel for analysis
func processBGPMessagesStatic(msgChannel chan []common.BGPMessage, windowChannel chan []common.Window) {
	var bucketMap = make(map[time.Time][]common.BGPMessage)

	for bgpMessages := range msgChannel {
		for _, msg := range bgpMessages {
			windowSize, _ := time.ParseDuration("60s")

			// Round down the timestamp to the nearest multiple of 60 seconds
			bucketTimestamp := msg.Timestamp.Truncate(windowSize)

			// Check if a bucket for the rounded timestamp exists, create it if not
			if _, ok := bucketMap[bucketTimestamp]; !ok {
				bucketMap[bucketTimestamp] = make([]common.BGPMessage, 0)
			}

			// Append the message to the corresponding bucket
			bucketMap[bucketTimestamp] = append(bucketMap[bucketTimestamp], msg)
		}
	}

	for timestamp, messages := range bucketMap {
		frequency := len(messages)
		fmt.Printf("Timestamp: %s, Frequency: %d\n", timestamp.Format("2006-01-02 15:04:05"), frequency)
	}
	fmt.Println("Minutes Analyzed: ", len(bucketMap))

	var windows []common.Window

	window := common.Window{
		Filter:    "none",
		BucketMap: bucketMap,
	}

	windows = append(windows, window)

	windowChannel <- windows

	close(windowChannel)

}
