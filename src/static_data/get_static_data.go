package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Folder struct {
	FolderName string
	Urls       []string
}

func downloadFolder(folder Folder) error {
	// Create the folder on the file system
	err := os.MkdirAll(folder.FolderName, os.ModePerm)
	if err != nil {
		return err
	}

	// Download files to the folder
	for _, url := range folder.Urls {
		err := downloadFile(url, folder.FolderName)
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadFile(url, folderDir string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	// Extract the filename from the URL
	_, filename := filepath.Split(url)

	// Create the file in the specified folder
	filePath := filepath.Join(folderDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the file from the response body to the local file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Downloaded: %s\n", filePath)
	return nil
}

func main() {
	folders := []Folder{}

	// Create folder instances
	bgptest1 := Folder{
		FolderName: "bgptest1",
		Urls: []string{
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211204.1845.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211204.1900.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211204.1915.bz2",
		},
	}
	folders = append(folders, bgptest1)

	bgptest2 := Folder{
		FolderName: "bgptest2",
		Urls: []string{
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211204.1545.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211204.1600.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211204.1615.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211204.1630.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211204.1645.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211204.1700.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211204.1715.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211204.1730.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211204.1745.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211204.1800.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211204.1815.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211204.1830.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211204.1845.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211204.1900.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211204.1915.bz2",
		},
	}
	folders = append(folders, bgptest2)

	bgptest3 := Folder{
		FolderName: "bgptest3",
		Urls: []string{
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211209.1530.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211209.1545.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211209.1600.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211209.1615.bz2",
		},
	}
	folders = append(folders, bgptest3)

	bgptest4 := Folder{
		FolderName: "bgptest4",
		Urls: []string{
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.1615.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.1630.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.1645.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.1700.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.1715.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.1730.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.1745.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.1800.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.1815.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.1830.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.1845.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.1900.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.1915.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.1930.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.1945.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.2000.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.2015.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.2030.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.2045.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.2100.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.2115.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.2130.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.2145.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.2200.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.2215.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.2230.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.2245.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.2300.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.2315.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211220.2330.bz2",
		},
	}
	folders = append(folders, bgptest4)

	bgptest5 := Folder{
		FolderName: "bgptest5",
		Urls: []string{
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.1800.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.1815.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.1830.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.1845.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.1900.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.1915.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.1930.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.1945.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.2000.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.2015.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.2030.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.2045.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.2100.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.2115.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.2130.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.2145.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.2200.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.2215.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.2230.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.2245.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.2300.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.2315.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.2330.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211228.2345.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211229.0000.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211229.0015.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211229.0030.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211229.0045.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211229.0100.bz2",
		},
	}
	folders = append(folders, bgptest5)

	bgpnyfiles := Folder{
		FolderName: "bgpnyfiles",
		Urls: []string{
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211201.0000.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211201.0015.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211201.0030.bz2",
			"http://routeviews.org/route-views.ny/bgpdata/2021.12/UPDATES/updates.20211201.0045.bz2",
		},
	}
	folders = append(folders, bgpnyfiles)

	bgpupdatefiles := Folder{
		FolderName: "bgpupdatefiles",
		Urls: []string{
			"http://routeviews.org/route-views.chicago/bgpdata/2016.06/UPDATES/updates.20160629.1330.bz2",
			"http://routeviews.org/route-views.chicago/bgpdata/2016.06/UPDATES/updates.20160629.1345.bz2",
			"http://routeviews.org/route-views.chicago/bgpdata/2016.06/UPDATES/updates.20160629.1400.bz2",
			"http://routeviews.org/route-views.chicago/bgpdata/2016.06/UPDATES/updates.20160629.1415.bz2",
			"http://routeviews.org/route-views.chicago/bgpdata/2016.06/UPDATES/updates.20160629.1430.bz2",
			"http://routeviews.org/route-views.chicago/bgpdata/2016.06/UPDATES/updates.20160629.1445.bz2",
		},
	}
	folders = append(folders, bgpupdatefiles)

	// Download folders
	for _, folder := range folders {
		err := downloadFolder(folder)
		if err != nil {
			fmt.Printf("Error downloading folder %s, %v", folder.FolderName, err)
		}
	}
}
