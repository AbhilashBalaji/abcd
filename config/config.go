package config

// Shard type from Config file
// describes shard which holds unique set of keys.
type Shard struct {
	Name    string
	Idx     int
	Address string
}

// Config  describes the sharding config
type Config struct {
	Shards []Shard
}
