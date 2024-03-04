package main

import (
	"fmt"
	"io"
	"log"
	"mira/config"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
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
	configStruct, err := config.LoadConfig("default-config.json")
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}
	err = configStruct.ValidateConfiguration()
	if err != nil {
		log.Fatalf("Failed to validate config: %v", err)
	}
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

	numFiles := configStruct.WindowSize
	partialUrls := extractBZ2URLs(doc)
	allFolders := createFolders(partialUrls, mainUrl, numFiles)
	for _, folder := range allFolders {
		err := downloadFolder(folder)
		if err != nil {
			return
		}
	}

}

func downloadFile(url, folderPath string) error {
	// Create the folder if it doesn't exist
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return fmt.Errorf("error creating folder: %v", err)
	}
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

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error downloading file, status code: %d", response.StatusCode)
	}

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

// this works
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
