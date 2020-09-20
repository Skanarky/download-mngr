package main

import (
	"fmt"
	"time"
)

type Download struct {
	OriginURI     string
	TargetPath    string
	TotalSections int
}

func main() {
	startTime := time.Now()
	// FIXME:
	// TODO: pass URI, create package?
	// TODO: create path dynamically 
	download := Download{
		OriginURI: "https://some-uri-will-gohere"
		TargetPath    "some-path-will-go-here.jpg"
		TotalSections 10
	}
	fmt.Printf("Executing download manger\n")
}
