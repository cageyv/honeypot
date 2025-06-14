package main

import (
	"flag"
	"os"
	"strconv"
	"testing"
)

func TestMainExitCode(t *testing.T) {
	tests := []struct {
		name     string
		envCode  string
		flagCode int
		want     int
	}{
		{"default", "", 0, 222},
		{"env_code", "111", 0, 111},
		{"flag_code", "", 201, 201},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags for each test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			exitCode = flag.Int("honeypot-exit-code", 222, "Exit code for honeypot (default: 222)")

			if tt.envCode != "" {
				os.Setenv("HONEYPOT_EXIT_CODE", tt.envCode)
				defer os.Unsetenv("HONEYPOT_EXIT_CODE")
			}

			if tt.flagCode != 0 {
				os.Args = []string{"cmd", "-honeypot-exit-code=" + strconv.Itoa(tt.flagCode)}
				flag.Parse()
			}

			code := getExitCode()
			if code != tt.want {
				t.Errorf("getExitCode() = %v, want %v", code, tt.want)
			}
		})
	}
}
