package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
)

const DEFAULT_PORT = "8080"

type config struct {
	Port string `json:"port" envconfig:"PORT"`

	// Can be file path or "-" for stdout.
	// "" or omitting the setting entirely disables logging
	LogFile string `json:"logFile" envconfig:"LOGFILE"`
}

func GetConfig() config {
	var cfg config
	var fname string

	// get config from commandline first, fallback to try environment
	if len(os.Args) > 2 {
		fname = os.Args[1]
	} else {
		fname = os.Getenv("CONFIG")
	}

	if fname != "" {
		f, err := os.Open(fname)
		if err == nil {
			defer f.Close()
			err := json.NewDecoder(f).Decode(&cfg)
			if err != nil {
				fmt.Printf("Warning: can't parse config file %v (ignoring): %v", fname, err.Error())
			}
		} else {
			fmt.Printf("Warning: can't open config file %v (ignoring): %v", fname, err.Error())
		}
	}

	envconfig.Process("", &cfg)

	// if all else fails, fallback to default
	if cfg.Port == "" {
		cfg.Port = DEFAULT_PORT
	}

	return cfg
}
