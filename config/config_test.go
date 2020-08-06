package config_test

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/AbhilashBalaji/abcd/config"
)

func CreateConfig(t *testing.T, contents string) config.Config {
	t.Helper()

	f, err := ioutil.TempFile(os.TempDir(), "config.toml")
	if err != nil {
		t.Fatalf("could not create temp file : %v", err)
	}
	defer f.Close()

	name := f.Name()
	defer os.Remove(name)

	_, err = f.WriteString(contents)

	if err != nil {
		t.Fatalf("could not write config contents : %v", err)
	}
	c, err := config.ParseFile(name)
	if err != nil {
		t.Fatalf("could not Parse config : %v", err)
	}
	return c

}

func TestConfigParse(t *testing.T) {
	got := CreateConfig(t, `[[shards]]
	name = "sh0"
	idx = 0
	address = "localhost:8080"`)
	want := config.Config{
		Shards: []config.Shard{
			{
				Name:    "sh0",
				Idx:     0,
				Address: "localhost:8080",
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Configs do not match got : %#v , want :%#v ", got, want)
	}
}

func TestParseShards(t *testing.T) {
	c := CreateConfig(t, `[[shards]]
	name = "sh0"
	idx = 0
	address = "localhost:8080"
	[[shards]]
	name = "sh1"
	idx = 1
	address = "localhost:8081"`)
	got, err := config.ParseShards(c.Shards, "sh1")
	if err != nil {
		t.Fatalf("could not parse shards %#v : %v", c.Shards, err)
	}

	want := &config.Shards{
		Count:  2,
		CurIdx: 1,
		Addrs: map[int]string{
			0: "localhost:8080",
			1: "localhost:8081",
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("The shards config do not match got : %#v , want :%#v ", got, want)
	}

}
