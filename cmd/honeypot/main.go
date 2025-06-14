package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

//go:embed binary_info.json
var embeddedBinaryInfo []byte

var (
	exitCode     = flag.Int("honeypot-exit-code", 222, "Exit code for honeypot (default: 222)")
	targetBinary = flag.String("target", "", "Target binary to emulate")
	analyzeMode  = flag.Bool("analyze", false, "Analyze target binary and print info")
	buildMode    = flag.Bool("build", false, "Build mode - analyze binaries and generate info file")
	replaceMode  = flag.Bool("replace", false, "Replace mode - create binary replacements")
)

// BinaryInfo is the information about a binary file
type BinaryInfo struct {
	Size     int64
	ExecTime time.Duration
	Header   []byte
	Path     string
}

var binaryCache map[string]BinaryInfo

func getExitCode() int {
	if envCode := os.Getenv("HONEYPOT_EXIT_CODE"); envCode != "" {
		if code, err := strconv.Atoi(envCode); err == nil {
			return code
		}
	}
	return *exitCode
}

func init() {
	if err := json.Unmarshal(embeddedBinaryInfo, &binaryCache); err != nil {
		log.Printf("Warning: Failed to load embedded binary info: %v", err)
		binaryCache = make(map[string]BinaryInfo)
	}
}

func findBinaryPath(name string) (string, error) {
	return exec.LookPath(name)
}

func analyzeBinary(path string) (BinaryInfo, error) {
	info := BinaryInfo{Path: path}

	file, err := os.Open(path)
	if err != nil {
		return info, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return info, err
	}
	info.Size = stat.Size()

	header := make([]byte, 16)
	_, err = io.ReadFull(file, header)
	if err != nil && err != io.EOF {
		return info, err
	}
	info.Header = header

	start := time.Now()
	cmd := exec.Command(path)
	cmd.Run()
	info.ExecTime = time.Since(start)

	return info, nil
}

func analyzeCommonBinaries() error {
	commonBins := []string{
		"ls", "cat", "sh", "sudo",
		"nc", "curl", "wget", "ssh",
		"python", "python3", "npm", "node",
		"apt", "apt-get", "yum", "dnf",
		"pip", "pip3", "git", "docker",
	}

	binaryCache = make(map[string]BinaryInfo)

	for _, bin := range commonBins {
		if path, err := findBinaryPath(bin); err == nil {
			info, err := analyzeBinary(path)
			if err != nil {
				log.Printf("Warning: Failed to analyze %s: %v", bin, err)
				continue
			}
			binaryCache[bin] = info
			log.Printf("Analyzed %s at %s", bin, path)
		} else {
			log.Printf("Warning: Binary %s not found in PATH", bin)
		}
	}

	data, err := json.MarshalIndent(binaryCache, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile("cmd/honeypot/binary_info.json", data, 0644); err != nil {
		return err
	}

	log.Printf("Generated binary info:\n%s", string(data))
	return nil
}

func emulateBinary() {
	if *targetBinary == "" {
		return
	}

	binaryName := filepath.Base(*targetBinary)
	info, exists := binaryCache[binaryName]

	if !exists {
		log.Printf("Warning: No info for binary %s", binaryName)
		return
	}

	// Emulate size
	dummy := make([]byte, info.Size)
	os.WriteFile("/tmp/honeypot_dummy", dummy, 0755)
	defer os.Remove("/tmp/honeypot_dummy")

	// Emulate execution time
	time.Sleep(info.ExecTime)
}

func createBinaryReplacements() error {
	commonBins := []string{
		"ls", "cat", "sh", "sudo",
		"nc", "curl", "wget", "ssh",
		"python", "python3", "npm", "node",
		"apt", "apt-get", "yum", "dnf",
		"pip", "pip3", "git", "docker",
	}

	honeypotPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get honeypot path: %v", err)
	}

	for _, bin := range commonBins {
		// Create binary copy
		if err := os.Link(honeypotPath, bin); err != nil {
			log.Printf("Warning: Failed to create replacement for %s: %v", bin, err)
			continue
		}
		log.Printf("Created replacement for %s", bin)
	}

	return nil
}

func main() {
	flag.Parse()

	// Check if we're running as a replacement binary
	execName := filepath.Base(os.Args[0])
	if execName != "honeypot" {
		*targetBinary = execName
	}

	if *buildMode {
		if err := analyzeCommonBinaries(); err != nil {
			log.Fatalf("Build failed: %v", err)
		}
		log.Println("Binary analysis complete. Generated binary_info.json")
		return
	}

	if *replaceMode {
		if err := createBinaryReplacements(); err != nil {
			log.Fatalf("Replace failed: %v", err)
		}
		log.Println("Binary replacements created")
		return
	}

	if *analyzeMode && *targetBinary != "" {
		info, err := analyzeBinary(*targetBinary)
		if err != nil {
			log.Fatalf("Analysis failed: %v", err)
		}
		log.Printf("Binary: %s\nSize: %d bytes\nExec Time: %v\nHeader: %x",
			*targetBinary, info.Size, info.ExecTime, info.Header)
		return
	}

	emulateBinary()
	os.Exit(getExitCode())
}
