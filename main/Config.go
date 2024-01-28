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

// Builds conf struct first from jsonConfig is it's not an empty string, then by checking if
// PORT or LOGFILE environment variables are set. NOTE: if a given setting is appears in both json and environment,
// environment value is used.
//
// Default Port = "8080"
// Default LogFile = "" (no logging).
//
// Set LogFile = "-" to log to stdout, otherwise provide
// a valid file name. If LogFile cannot be accessed and parsed, a warning is printed and
// execution continues with either environment values or defaults.
func GetConfig(jsonConfig string) config {
	var cfg config

	if jsonConfig != "" {
		f, err := os.Open(jsonConfig)
		if err == nil {
			defer f.Close()
			err := json.NewDecoder(f).Decode(&cfg)
			if err != nil {
				fmt.Printf("Warning: can't parse config file %v (ignoring): %v", jsonConfig, err.Error())
			}
		} else {
			fmt.Printf("Warning: can't open config file %v (ignoring): %v", jsonConfig, err.Error())
		}
	}

	envconfig.Process("", &cfg)

	// if all else fails, fallback to default
	if cfg.Port == "" {
		cfg.Port = DEFAULT_PORT
	}

	return cfg
}
