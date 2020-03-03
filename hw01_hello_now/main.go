package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

const urlNtp = "0.beevik-ntp.pool.ntp.org"

func main() {
	currentTime := time.Now()
	exactTime, err := ntp.Time(urlNtp)
	if err != nil {
		log.Fatalf("ERROR get time from %s: %s", urlNtp, err)
	}
	fmt.Printf("current time: %v\n", currentTime)
	fmt.Printf("exact time: %v\n", exactTime)
}
