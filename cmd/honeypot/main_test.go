package main

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestMain(m *testing.M) {
	// Create empty binary_info.json for tests
	emptyCache := make(map[string]BinaryInfo)
	data, _ := json.Marshal(emptyCache)
	os.WriteFile("binary_info.json", data, 0644)
	defer os.Remove("binary_info.json")

	os.Exit(m.Run())
}

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

func TestBinaryReplacement(t *testing.T) {
	// Create temp dir for test
	tmpDir, err := os.MkdirTemp("", "honeypot-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp dir
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer os.Chdir(oldDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Create test honeypot binary
	testBin := filepath.Join(tmpDir, "honeypot")
	if err := os.WriteFile(testBin, []byte("test"), 0755); err != nil {
		t.Fatalf("Failed to create test binary: %v", err)
	}

	// Test replacement creation
	os.Args = []string{"honeypot", "-replace"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	replaceMode = flag.Bool("replace", false, "Replace mode")
	flag.Parse()

	if err := createBinaryReplacements(); err != nil {
		t.Errorf("createBinaryReplacements() error = %v", err)
	}

	// Verify replacements exist
	for _, bin := range []string{"ls", "cat", "sh"} {
		if _, err := os.Stat(bin); err != nil {
			t.Errorf("Replacement for %s not created: %v", bin, err)
		}
	}
}

func TestBinaryEmulation(t *testing.T) {
	// Test running as replacement binary
	os.Args = []string{"cat"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	targetBinary = flag.String("target", "cat", "Target binary")
	flag.Parse()

	// Verify target is set correctly
	if *targetBinary != "cat" {
		t.Errorf("targetBinary = %v, want cat", *targetBinary)
	}
}
