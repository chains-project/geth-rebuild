package experiments

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/chains-project/geth-rebuild/internal/utils"
)

// utilities to generate experiment inputs, i.e. inputs to the main program in gethrebuild

type ExperimentInput struct {
	OS       utils.OS
	Arch     utils.Arch
	Version  string
	Unstable string
}

func (i ExperimentInput) String() string {
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

// Generates all experiment inputs for given os
func GenerateAllExperiments(ops utils.OS, arches []utils.Arch, stableVersions []string, commitsFile string) (experiments [][]ExperimentInput, err error) {
	for _, arch := range arches {
		// Generate stable experiments
		exps, err := GenerateExperimentInputs(utils.Linux, arch, stableVersions, "")
		if err != nil {
			return experiments, fmt.Errorf("error generating stable version experiments: %v", err)
		}
		experiments = append(experiments, exps)

		// Generate unstable experiments
		exps, err = GenerateExperimentInputs(utils.Linux, arch, nil, commitsFile)
		if err != nil {
			return experiments, fmt.Errorf("error generating unstable version experiments: %v", err)
		}
		experiments = append(experiments, exps)
	}
	return experiments, nil
}

// Generates all experiments for given os/arch
func GenerateExperimentInputs(ops utils.OS, arch utils.Arch, stableVersions []string, commitsFile string) (experiments []ExperimentInput, err error) {
	if commitsFile != "" {
		data, err := os.ReadFile(commitsFile)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}

		var jsonData JSONData
		if err := json.Unmarshal(data, &jsonData); err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
		}

		for _, commit := range jsonData.Commits {
			version := strings.TrimPrefix(commit.Version, "v")
			experiments = append(experiments, ExperimentInput{OS: ops, Arch: arch, Version: version, Unstable: commit.SHA})
		}
	} else {
		// Generate stable experiments if no filepath is provided
		for _, version := range stableVersions {
			experiments = append(experiments, ExperimentInput{OS: ops, Arch: arch, Version: version, Unstable: ""})
		}
	}

	return experiments, nil
}
