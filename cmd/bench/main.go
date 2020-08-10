package main

import (
	"flag"
	"fmt"
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
	addr        = flag.String("addr", "localhost:8080", "HTTP host port for benchmarking")
	iterations  = flag.Int("iterations", 1000, "Number of iterations")
	concurrency = flag.Int("concurency", 4, "Level of concurrency")
)

func benchmark(name string, fn func()) {
	var max time.Duration
	var min = time.Hour
	start := time.Now()
	for i := 0; i < *iterations; i++ {
		iterStart := time.Now()
		fn()
		iterTime := time.Since(iterStart)

		if iterTime > max {
			max = iterTime
		}
		if iterTime < min {
			min = iterTime
		}
	}
	qps := float64(N) / (float64(time.Since(start)) / float64(time.Second))
	fmt.Printf("Func %s took %s avg, %.1f QPS, %s max, %s min\n", name, time.Since(start)/N, qps, max, min)
}

func writeRand() {
	key := fmt.Sprintf("key-%d", rand.Intn(10000))
	value := fmt.Sprintf("value-%d\n", rand.Intn(10000))

	values := url.Values{}

	values.Set("key", key)
	values.Set("value", value)

	resp, err := http.Get("http://" + (*addr) + "/set?" + values.Encode())
	if err != nil {
		log.Fatalf("Error during set :%v", err)
	}
	defer resp.Body.Close()

	// fmt.Printf("key = %s value = %s\n", key, value)
}
func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	fmt.Printf("Running with %d iterations and concurrency level %d\n", *iterations, *concurrency)

	var wg sync.WaitGroup

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			benchmark("write", writeRand)
			wg.Done()
		}()
	}

	wg.Wait()

}
