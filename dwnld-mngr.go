package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Download struct {
	OriginURI     string
	TargetPath    string
	TotalSections int
}

func main() {
	startTime := time.Now()

	download := Download{
		OriginURI:     "https://www.pexels.com/photo/573271/download/",
		TargetPath:    "the-pic.jpg",
		TotalSections: 10,
	}

	err := download.Run()

	checkError(err)

	fmt.Printf("Download completed in: %v\n", time.Now().Sub(startTime).Seconds())
}

func (download Download) Run() error {
	fmt.Println("Download started. Connecting...")

	req, err := download.PrepareNewRequest("HEAD")
	checkError(err)

	res, err := http.DefaultClient.Do(req)
	checkError(err)
	if res.StatusCode > 299 {
		return errors.New(fmt.Sprintf("Download can not be finished. Status %v", res.StatusCode))
	}

	size, err := strconv.Atoi(res.Header.Get("Content-Length"))
	checkError(err)
	fmt.Printf("Size is %v bytes\n", size)

	sections := make([][2]int, download.TotalSections)
	oneSectSize := size / download.TotalSections

	for i := range sections {
		if i == 0 {
			sections[i][0] = 0
		} else {
			sections[i][0] = sections[i-1][1] + 1
		}
		if i < download.TotalSections-1 {
			sections[i][1] = sections[i][0] + oneSectSize
		} else {
			sections[i][1] = size - 1
		}
	}

	for i, sect := range sections {
		err := download.DownloadOneSection(i, sect)
		checkError(err)
	}

	return nil
}

func (download Download) DownloadOneSection(i int, sect [2]int) error {
	req, err := download.PrepareNewRequest("GET")
	checkError(err)

	req.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", sect[0], sect[1]))
	res, err := http.DefaultClient.Do(req)
	checkError(err)
	fmt.Printf("Dowloaded %v bytes for section No. %v: %v\n", res.Header.Get("Content-Length"), i, sect)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	checkError(err)

	err = ioutil.WriteFile(fmt.Sprintf("sect-%v.tmp", i), bodyBytes, os.ModePerm)
	checkError(err)
	return nil
}
func (download Download) PrepareNewRequest(method string) (*http.Request, error) {
	req, err := http.NewRequest(
		method,
		download.OriginURI,
		nil,
	)
	checkError(err)
	req.Header.Set("User-Agent", "Simple Download Manager")
	return req, nil
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Error while downloading: %s\n", err)
	}
}
