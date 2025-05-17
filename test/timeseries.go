package main

import (
	"fmt"
	"log_detector/timeseries"
	"runtime"
	"time"
)

func main() {
	ts, err := timeseries.NewTimeSeries()
	if err != nil {
		// handle error
	}

	blockChan := make(chan string)
	go func() {
		for i := range 100000 {
			time.Sleep(1 * time.Second)
			fmt.Println("increase: ", i)
			ts.Increase(i)
		}
		blockChan <- "Done"
	}()

	go func() {
		for {
			rangeValue, err := ts.Range(time.Now().Add(-2*time.Second), time.Now())
			recentValue, err := ts.Recent(60 * time.Second)
			if err != nil {
				fmt.Printf("%v", err)
				return
			}
			fmt.Printf("time series result : %f, %f \n", rangeValue, recentValue)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	go func() {
		for {
			printMemUsage()
			time.Sleep(2000 * time.Millisecond)
		}
	}()

	<-blockChan

}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v KiB", m.Alloc/1024)
	fmt.Printf("\tTotalAlloc = %v KiB", m.TotalAlloc/1024)
	fmt.Printf("\tSys = %v KiB", m.Sys/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
