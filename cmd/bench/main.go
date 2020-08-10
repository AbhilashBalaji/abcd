package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// N is number of iterations in benchmark
const N = 1000

var (
	addr           = flag.String("addr", "127.0.0.2:8080", "HTTP host port for benchmarking")
	iterations     = flag.Int("iterations", 10000, "Number of iterations")
	readiterations = flag.Int("read-iterations", 1000, "Number of read iterations")
	concurrency    = flag.Int("concurency", 2, "Level of concurrency")
)

var httpClient = &http.Client{
	Transport: &http.Transport{
		IdleConnTimeout:     time.Second * 60,
		MaxIdleConns:        32,
		MaxConnsPerHost:     32,
		MaxIdleConnsPerHost: 32,
	},
}

func benchmark(name string, iterations int, fn func() string) (qps float64, strs []string) {
	var max time.Duration
	var min = time.Hour
	start := time.Now()
	for i := 0; i < iterations; i++ {
		iterStart := time.Now()
		// fn()
		strs = append(strs, fn())
		iterTime := time.Since(iterStart)

		if iterTime > max {
			max = iterTime
		}
		if iterTime < min {
			min = iterTime
		}
	}
	avg := time.Since(start) / N
	qps = float64(N) / (float64(time.Since(start)) / float64(time.Second))
	fmt.Printf("Func %s took %s avg, %.1f QPS, %s max, %s min\n", name, avg, qps, max, min)

	return qps, strs
}

func writeRand() (key string) {
	key = fmt.Sprintf("key-%d", rand.Intn(10000))
	value := fmt.Sprintf("value-%d\n", rand.Intn(10000))

	values := url.Values{}

	values.Set("key", key)
	values.Set("value", value)

	resp, err := httpClient.Get("http://" + (*addr) + "/set?" + values.Encode())
	if err != nil {
		log.Fatalf("Error during set :%v", err)
	}
	defer resp.Body.Close()
	return key

}

func readRand(allKeys []string) (key string) {
	key = allKeys[rand.Intn(len(allKeys))]

	values := url.Values{}

	values.Set("key", key)

	resp, err := httpClient.Get("http://" + (*addr) + "/get?" + values.Encode())
	if err != nil {
		log.Fatalf("Error during get :%v", err)
	}
	// Keep alive  only works when read and  wont break benchmarking
	io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	return key
}

func benchmarkWrite() (allKeys []string) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var totalQPS float64

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			qps, strs := benchmark("write", *iterations, writeRand)
			mu.Lock()
			totalQPS += qps
			allKeys = append(allKeys, strs...)
			mu.Unlock()
			wg.Done()
		}()
	}

	wg.Wait()
	log.Printf("WRITE Total QPS : %.1f , set %d keys\n", totalQPS, len(allKeys))
	return allKeys

}

func benchmarkRead(allKeys []string) {

	var totalQPS float64
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			qps, _ := benchmark("write", *readiterations, func() string { return readRand(allKeys) })
			mu.Lock()
			totalQPS += qps

			mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()

	log.Printf("READ Total QPS : %.1f \n", totalQPS)

}

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	fmt.Printf("Running with %d iterations and concurrency level %d\n", *iterations, *concurrency)

	allKeys := benchmarkWrite()
	go benchmarkWrite()
	benchmarkRead(allKeys)

}
