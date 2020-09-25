package main

import (
	"fmt"
	"log"
	"net/http"
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
		// OriginURI:     "https://drive.google.com/uc?id=16-ik4vC6ynIfQP3_DdziQH9KyjeIdVlM&authuser=0&export=download",
		OriginURI:     "https://drive.google.com/file/d/16-ik4vC6ynIfQP3_DdziQH9KyjeIdVlM/view?usp=sharing",
		TargetPath:    "the-pic.jpg",
		TotalSections: 10,
	}

	err := download.Run()

	checkError(err)

	fmt.Printf("Download completed in: %v\n", time.Now().Sub(startTime).Seconds())
}

func (download Download) Run() error {
	fmt.Println("Download started. Connecting.")
	req, err := download.PrepareNewRequest("HEAD")
	checkError(err)
	res, err := http.DefaultClient.Do(req)
	checkError(err)
	fmt.Printf("Response Stat Code: %v\n", res.StatusCode)
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
