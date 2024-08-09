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

// All stable builds since 1.14.2, which is the upgrade from ubuntu bionic to noble in travis.yml
var stableVersions = []string{"1.14.7", "1.14.6", "1.14.5", "1.14.4", "1.14.3", "1.14.2"}

func GenerateExperimentInputs(ops utils.OS, arch utils.Arch, commitsFile string) (experiments []ExperimentInput, err error) {
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
