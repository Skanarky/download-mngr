package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
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

	if err != nil {
		log.Fatalf("Error while downloading: %s\n", err)
	}

	fmt.Printf("Download completed in: %v\n", time.Now().Sub(startTime).Seconds())
}

func (download Download) Run() error {
	fmt.Println("Download started. Connecting...")

	req, err := download.PrepareNewRequest("HEAD")
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode > 299 {
		return errors.New(fmt.Sprintf("Download can not be finished. Status %v", res.StatusCode))
	}

	size, err := strconv.Atoi(res.Header.Get("Content-Length"))
	if err != nil {
		return err
	}
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

	var waitGr sync.WaitGroup
	for i, sect := range sections {
		waitGr.Add(1)
		inI, inSect := i, sect
		go func() {
			defer waitGr.Done()
			err := download.DownloadOneSection(inI, inSect)
			if err != nil {
				panic(err)
			}
		}()
	}
	waitGr.Wait()

	err = download.MergeAllSectFiles(sections)
	if err != nil {
		return err
	}

	return nil
}

func (download Download) DownloadOneSection(i int, sect [2]int) error {
	req, err := download.PrepareNewRequest("GET")
	if err != nil {
		return err
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", sect[0], sect[1]))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	fmt.Printf("Dowloaded %v bytes for section No. %v: %v\n", res.Header.Get("Content-Length"), i, sect)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("sect-%v.tmp", i), bodyBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
func (download Download) PrepareNewRequest(method string) (*http.Request, error) {
	req, err := http.NewRequest(
		method,
		download.OriginURI,
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Simple Download Manager")
	return req, nil
}
func (download Download) MergeAllSectFiles(sections [][2]int) error {
	theFile, err := os.OpenFile(download.TargetPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}

	defer theFile.Close()

	for i := range sections {
		tempBytes, err := ioutil.ReadFile(fmt.Sprintf("sect-%v.tmp", i))
		if err != nil {
			return err
		}

		writtenBytes, err := theFile.Write(tempBytes)
		if err != nil {
			return err
		}

		fmt.Printf("%v bytes were merged\n", writtenBytes)
	}

	return nil
}
