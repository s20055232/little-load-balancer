package main

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type MyConfig struct {
	Version int
	Title   string
	Host    string
	Port    int
	Servers []string
}

func readTOMLSetting() MyConfig {
	var cfg MyConfig
	doc, err := os.ReadFile("config.toml")
	check(err)
	err = toml.Unmarshal([]byte(doc), &cfg)
	check(err)
	return cfg
}
