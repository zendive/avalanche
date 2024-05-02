package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/zendive/avalanche/lib"
)

var testUrl string = "http://localhost:8282/"
var testStartTime = time.Now()
var meanResponse = lib.NewAverageble(lib.AveragebleOptions{Population: true})
var meanFetchesPerSec = lib.NewAverageble(lib.AveragebleOptions{Population: true})
var meanSuccess = lib.NewAverageble(lib.AveragebleOptions{Population: true})
var fetchesNum uint64 = 0
var successNum uint64 = 0
var mutex = &sync.Mutex{}
var responseStatusCount = make(map[interface{}]int64)
var secElapsed = time.Duration(0)
var MAX_RUNTIME = time.Duration(10) * time.Minute

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Printf("Commencing avalanche fetch on: %s\n", testUrl)
	fmt.Println("Press Ctrl+C to stop...")

	ticker := time.NewTicker(time.Duration(1) * time.Second)
	chanSysSignals := make(chan os.Signal, 2)
	signal.Notify(chanSysSignals, os.Interrupt, syscall.SIGTERM)

	// whait for Ctrl+C
	go func() {
		<-chanSysSignals
		pprof.StopCPUProfile()
		ticker.Stop()
		printStats()
		os.Exit(0)
	}()

	// each tick: update fps
	go func() {
		for range ticker.C {
			fmt.Print(successNum, " ")
			meanFetchesPerSec.Add(float64(fetchesNum))
			meanSuccess.Add(float64(successNum))
			atomic.AddUint64(&fetchesNum, -fetchesNum)
			atomic.AddUint64(&successNum, -successNum)

			secElapsed += time.Duration(1) * time.Second
			if secElapsed >= MAX_RUNTIME {
				chanSysSignals <- syscall.SIGTERM
			}
		}
	}()

	client := http.Client{}

	// start fetching
	for i := 0; ; i++ {
		url := fmt.Sprintf("%s?avalanche=%d", testUrl, i)
		go fetchUrl(&url, &client)
	}
}

func fetchUrl(url *string, client *http.Client) {
	atomic.AddUint64(&fetchesNum, 1)

	startTime := time.Now()
	response, err := client.Get(*url)

	if err != nil {
		//fmt.Println(err)
		mutex.Lock()
		responseStatusCount["total"] += 1
		responseStatusCount["error"] += 1
		mutex.Unlock()
	} else {
		defer response.Body.Close()

		mutex.Lock()
		if http.StatusOK == response.StatusCode {
			atomic.AddUint64(&successNum, 1)
		}
		responseStatusCount["total"] += 1
		responseStatusCount[fmt.Sprintf("%v", response.StatusCode)] += 1
		meanResponse.Add(time.Since(startTime).Seconds())
		mutex.Unlock()
	}
}

func printStats() {
	mutex.Lock()
	fmt.Printf("\nTest complete in:\t%v for url: %s\n", time.Since(testStartTime), testUrl)
	fmt.Printf("Successfull (rps):\t%v\n", meanSuccess)
	fmt.Printf("Response time (s):\t%v\n", meanResponse)
	fmt.Printf("Fetch rate (fps):\t%v\n", meanFetchesPerSec)
	fmt.Printf("Σ of status codes:\t%v\n", responseStatusCount)
	fmt.Println("---")
	mutex.Unlock()
}
