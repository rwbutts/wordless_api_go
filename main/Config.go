package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
)

const CFG_JSON_FILE = "./config.json"
const DEFAULT_LISTEN_ADDRESS = ":80"

type config struct {
	ListenAddress string `json:"listenAddress" envconfig:"LISTEN_ADDRESS"`

	// Can be file path or "-" for stdout.
	// "" or omitting the setting entirely disables logging
	LogFile string `json:"logFile"`
}

func GetConfig() config {
	var cfg config

	f, err := os.Open(CFG_JSON_FILE)
	if err == nil {
		defer f.Close()
		err := json.NewDecoder(f).Decode(&cfg)
		if err != nil {
			fmt.Println("Warning: can't parse config file (ignoring): " + err.Error())
		}
	}

	envconfig.Process("", &cfg)

	// if all else fails, fallback to default
	if cfg.ListenAddress == "" {
		cfg.ListenAddress = DEFAULT_LISTEN_ADDRESS
	}

	return cfg
}
