package config

import (
	"fmt"
	"hash/fnv"

	"github.com/BurntSushi/toml"
)

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

// Shards represents an easier-to-use representation of
// the sharding config: the shards count, current index and
// the addresses of all other shards too
type Shards struct {
	Count  int
	CurIdx int
	Addrs  map[int]string
}

// ParseFile parses config and returns it on success.
func ParseFile(filename string) (Config, error) {
	var c Config

	if _, err := toml.DecodeFile(filename, &c); err != nil {
		return Config{}, err
	}
	return c, nil
}

// ParseShards converts and verifies the list of shards
// specified in config into a form
func ParseShards(shards []Shard, curShardName string) (*Shards, error) {
	shardCount := len(shards)
	shardIdx := -1
	addrs := make(map[int]string)

	for _, s := range shards {
		if _, ok := addrs[s.Idx]; ok {
			return nil, fmt.Errorf("duplicate shard idx : %d", s.Idx)
		}
		addrs[s.Idx] = s.Address
		if s.Name == curShardName {
			shardIdx = s.Idx
		}
	}

	for i := 0; i < shardCount; i++ {
		if _, ok := addrs[i]; !ok {
			return nil, fmt.Errorf("Shard %d was not fount", i)
		}
	}
	if shardIdx < 0 {
		return nil, fmt.Errorf("Shard %q was not found", curShardName)
	}
	return &Shards{
		Addrs:  addrs,
		Count:  shardCount,
		CurIdx: shardIdx,
	}, nil

}

// Index returns shard number for given key
func (s *Shards) Index(key string) int {
	h := fnv.New64()
	h.Write([]byte(key))
	return int(h.Sum64() % uint64(s.Count))
}
