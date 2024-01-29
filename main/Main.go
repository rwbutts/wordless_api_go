package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
)

const EXENAME = "wordless"
const VERSION = "1.7.4"
const HTTP_VER_HEADER = "X-wordless-api-version"

func main() {
	helpSwitches := []string{"--help", "-h", "/?"}

	if len(os.Args) > 1 && slices.Contains(helpSwitches, os.Args[1]) {
		fmt.Printf("Usage: %v (%v) [json config file]\n", EXENAME, VERSION)
		fmt.Println("\nSample file:")
		fmt.Println("    {")
		fmt.Println(`		  "port" : "8080"`)
		fmt.Println(`		  "logfile" : "config.json"`)
		fmt.Println("    }")
		fmt.Println("\nValues can also be set by the environment variables: PORT and LOGFILE.")
		fmt.Println("Note: Environment variable has higher priority than config file")
		fmt.Println("\nDefault port is 8080; default logfile = \"\" (no logging)")
		fmt.Println(`To log to stdout, set logfile =  "-".`)

		os.Exit(0)
	}

	var configFile string
	// get config from commandline first, fallback to try environment
	if len(os.Args) > 2 {
		configFile = os.Args[1]
	} else {
		// if not present, the final result is "" -- no config file is loaded
		configFile = os.Getenv("CONFIG")
	}

	config := GetConfig(configFile)

	setLogDestination(config.LogFile)

	fmt.Println("HTTP listening on :" + config.Port)
	log.Println("HTTP listening on :" + config.Port)

	mux := http.NewServeMux()
	MapRoutes(mux)

	mw := LoggingMiddleware(log.Default())

	err := http.ListenAndServe(":"+config.Port, mw(mux))
	if err != nil {
		log.Println("PANIC ListenAndServe() exit: " + err.Error())
		panic(err.Error())
	}
}

// "" = no log, "-" = stdout, otherwise a file name
// If file givem and cannot be opened, a warning is printed, but execution continues.
func setLogDestination(logOutput string) {
	if logOutput == "" {
		log.SetOutput(io.Discard)
	} else if logOutput == "-" {
		log.SetOutput(os.Stdout)
	} else {
		logfile, err := os.OpenFile(logOutput, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Printf("Warning: error opening log file %v: %v", logOutput, err.Error())
		} else {
			log.SetOutput(logfile)
		}
	}

}
