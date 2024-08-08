package experiments

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/chains-project/geth-rebuild/internal/utils"
)

type Input struct {
	OS       utils.OS
	Arch     utils.Arch
	Version  string
	Unstable string
}

func (i Input) String() string {
	return fmt.Sprintf("\n[EXPERIMENT INPUT] OS: %s, Arch: %s, Version: %s, Unstable: %s",
		i.OS, i.Arch, i.Version, i.Unstable)
}

type Commit struct {
	SHA     string `json:"commit"`
	Version string `json:"version"`
}

type JSONData struct {
	Since   string   `json:"since"`
	To      string   `json:"to"`
	Commits []Commit `json:"commits"`
}

// All stable builds since 1.14.2, which is the upgrade from ubuntu bionic to noble in travis.yml
var versions []string = []string{"1.14.7", "1.14.6", "1.14.5", "1.14.4", "1.14.3", "1.14.2"}

func generateStableExperiments(ops utils.OS, arch utils.Arch) (experiments []Input) {
	for _, version := range versions {
		experiments = append(experiments, Input{OS: ops, Arch: arch, Version: version, Unstable: ""})
	}
	return experiments
}

func generateUnstableExperiments(ops utils.OS, arch utils.Arch, filepath string) (experiments []Input) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var jsonData JSONData
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	for _, commit := range jsonData.Commits {
		experiments = append(experiments, Input{OS: ops, Arch: arch, Version: commit.Version, Unstable: commit.SHA})
	}
	return experiments
}

func GenerateAllExperiments() {
	var stables [][]Input
	var unstables [][]Input

	base, err := utils.GetRootDir("geth-rebuild")
	if err != nil {
		log.Fatal(fmt.Errorf("Error finding root path: %w", err))
	}
	path := filepath.Join(base, "experiments", "random_commits.json")

	arches := []utils.Arch{utils.AMD64, utils.A386, utils.ARM5, utils.ARM6, utils.ARM7, utils.ARM64}

	for _, arch := range arches {
		stables = append(stables, generateStableExperiments(utils.Linux, arch))
		unstables = append(unstables, generateUnstableExperiments(utils.Linux, utils.AMD64, path))

	}
	fmt.Printf("\n[STABLES]:\n\n %v", stables)
	fmt.Printf("\n[UNSTABLES]:\n\n %v", unstables)
}
