## distributed KV 

* Example Launch config provided  ```ExampleLaunch.sh```
* ```cmd/bench.go``` has a benchmarking tool
* ```sharding.toml``` contains sharding config

# Benchmarks

Running with 10000 iterations and concurrency level 2 \
Func write took 37.430405ms avg, 26.7 QPS, 16.703904ms max, 1.89063ms min\
WRITE Total QPS : 53.4 , set 20000 keys\
Func write took 224.856µs avg, 4447.3 QPS, 934.657µs max, 78.68µs min\
READ Total QPS : 8862.1 