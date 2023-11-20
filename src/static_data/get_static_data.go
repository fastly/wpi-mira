package main

import (
	"BGPAlert/analyze"
	"BGPAlert/blt_mad"
	"BGPAlert/common"
	"BGPAlert/config"
	"BGPAlert/optimization"
	"BGPAlert/parse"
	"BGPAlert/process"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
			sortedFrequencies := analyze.GetSortedFrequencies(lengthMap)

			// Convert to float array for analysis functions
			floatArray := analyze.Int64ArrayToFloat64Array(sortedFrequencies)
			//save frequency array into the freqOutFile
			blt_mad.SaveArrayToFile(freqOutFile, floatArray)

			//get min reqArray for the 97th percentile
			minReqArray := blt_mad.GetValuesLargerThanPercentile(floatArray, 97)
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
		parse.ParseStaticFile(inputFolderPath, msgChannel)
		wg.Done()
	}()

	go func() {
		process.ProcessBGPMessages(msgChannel, windowChannel)
		wg.Done()
	}()

	go func() {
		AnalyzeBGPMessagesWriteOntoFile(windowChannel, outFile, outMinReqFile)
		wg.Done()
	}()

	wg.Wait()

}
