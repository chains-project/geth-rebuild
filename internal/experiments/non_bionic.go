package experiments

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/chains-project/geth-rebuild/internal/buildconfig"
	"github.com/chains-project/geth-rebuild/internal/utils"
)

// Define a struct to match the JSON structure
type CommitInfo struct {
	Commit  string `json:"commit"`
	Version string `json:"version"`
}

// Define a struct for logging non-bionic distributions
type DistInfoLog struct {
	Commit  string `json:"commit"`
	Version string `json:"version"`
	Dist    string `json:"dist"`
}

func GetAvailableUnstable(scriptPath string, commit1 string, commit2 string, v1 string, v2 string, gethDir string) {
	_, err := utils.RunCommand(scriptPath, commit1, commit2, v1, v2, gethDir)
	if err != nil {
		log.Fatalf("Failed to get unstable commits: %v", err)
	}
}

func FindDistForCommits(commitsFile string, writeTo string) {
	// Open and read the JSON file
	file, err := os.Open(commitsFile)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Read the file's content
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Parse the JSON data
	var commits []CommitInfo
	err = json.Unmarshal(data, &commits)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Slice to store logs of non-bionic distributions
	var distLog []DistInfoLog

	// Iterate through each commit and call the function
	for _, commit := range commits {
		ubuntuDist, err := buildconfig.GetUbuntuDist(commit.Commit, 8)
		if err != nil {
			log.Printf("Error getting Ubuntu distribution for commit %s: %v", commit.Commit, err)
			continue
		}
		if ubuntuDist != "" {
			logEntry := DistInfoLog{
				Commit:  commit.Commit,
				Version: commit.Version,
				Dist:    ubuntuDist,
			}
			distLog = append(distLog, logEntry)
		}
	}

	logFile, err := os.Create(writeTo)
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()

	logData, err := json.MarshalIndent(distLog, "", "    ")
	if err != nil {
		log.Fatalf("Failed to marshal log data: %v", err)
	}

	_, err = logFile.Write(logData)
	if err != nil {
		log.Fatalf("Failed to write log data: %v", err)
	}

	fmt.Printf("\nFinished logging to %s\n", writeTo)
}
