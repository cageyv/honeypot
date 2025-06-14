package main

import (
	"flag"
	"os"
	"strconv"
)

var exitCode = flag.Int("honeypot-exit-code", 222, "Exit code for honeypot (default: 222)")

func getExitCode() int {
	if envCode := os.Getenv("HONEYPOT_EXIT_CODE"); envCode != "" {
		if code, err := strconv.Atoi(envCode); err == nil {
			return code
		}
	}
	return *exitCode
}

func main() {
	flag.Parse()
	os.Exit(getExitCode())
}
