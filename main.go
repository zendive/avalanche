package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/zendive/avalanche/lib"
)

var testUrl string = "http://localhost:8282/"
var testStartTime = time.Now().UTC()
var meanResponse = lib.NewAverageble(lib.AveragebleOptions{Population: true})
var meanFetchesPerSec = lib.NewAverageble(lib.AveragebleOptions{Population: true})
var meanSuccess = lib.NewAverageble(lib.AveragebleOptions{Population: true})
var fetchesNum uint64 = 0
var successNum uint64 = 0
var mutex = &sync.Mutex{}
var responseStatusCount = make(map[interface{}]int64)
var MAX_RUNTIME_DURATION = time.Duration(1) * time.Minute

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if len(os.Args) > 1 && os.Args[1] != "" {
		testUrl = os.Args[1]
	}

	fmt.Printf("Version: %s, CPUs: %d\n", runtime.Version(), runtime.NumCPU())
	fmt.Println("Press Ctrl+C to stop...")

	ticker := time.NewTicker(time.Duration(1) * time.Second)
	chanSysSignals := make(chan os.Signal, 1)
	signal.Notify(chanSysSignals, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		<-chanSysSignals
		ticker.Stop()
		printStats()
		os.Exit(0)
	}()

	// each tick: update fps
	go func() {
		for range ticker.C {
			fmt.Printf("%d ", successNum)
			meanFetchesPerSec.Add(float64(fetchesNum))
			meanSuccess.Add(float64(successNum))
			atomic.AddUint64(&fetchesNum, -fetchesNum)
			atomic.AddUint64(&successNum, -successNum)

			elapsed := time.Now().UTC().Sub(testStartTime)
			if elapsed >= MAX_RUNTIME_DURATION {
				chanSysSignals <- syscall.SIGINT
			}
		}
	}()

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := http.Client{Transport: customTransport}

asyncFetching:
	for i := 0; ; i++ {
		select {
		// whait for Ctrl+C
		case <-chanSysSignals:
			break asyncFetching

		default:
			url := fmt.Sprintf("%s?avalanche=%d", testUrl, i)
			go fetchUrl(&url, &client)
			// shortening sleep duration generates more errors on tested server,
			// but also suffocate current process
			time.Sleep(time.Millisecond)
		}
	}
}

func fetchUrl(url *string, client *http.Client) {
	atomic.AddUint64(&fetchesNum, 1)

	startTime := time.Now()
	response, err := client.Get(*url)

	if err != nil {
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
	fmt.Printf("\nTest complete in:\t%.1f(s) for url: %s\n", time.Since(testStartTime).Seconds(), testUrl)
	fmt.Printf("Successfull (rps):\t%v\n", meanSuccess)
	fmt.Printf("Response time (s):\t%v\n", meanResponse)
	fmt.Printf("Fetch rate (fps):\t%v\n", meanFetchesPerSec)
	fmt.Printf("Σ of status codes:\t%v\n", responseStatusCount)
	fmt.Println("---")
	mutex.Unlock()
}
