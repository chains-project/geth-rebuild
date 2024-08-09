package main

import (
	"log"
	"path/filepath"

	"github.com/chains-project/geth-rebuild/experiments"
	"github.com/chains-project/geth-rebuild/internal/utils"
)

func main() {
	var stableExps [][]experiments.ExperimentInput
	var unstableExps [][]experiments.ExperimentInput

	base, err := utils.GetRootDir("geth-rebuild")
	if err != nil {
		log.Fatalf("error finding root path: %v", err)
	}
	executable := filepath.Join(base, "gethrebuild")

	path := filepath.Join(base, "experiments", "random_commits.json")
	arches := []utils.Arch{utils.AMD64, utils.A386, utils.ARM5, utils.ARM6, utils.ARM7, utils.ARM64}

	// Generate experiments for each architecture
	for _, arch := range arches {
		// Generate stable experiments
		exps, err := experiments.GenerateExperimentInputs(utils.Linux, arch, "")
		if err != nil {
			log.Fatalf("error generating stable version experiments: %v", err)
		}
		stableExps = append(stableExps, exps)

		// Generate unstable experiments
		exps, err = experiments.GenerateExperimentInputs(utils.Linux, arch, path)
		if err != nil {
			log.Fatalf("error generating unstable version experiments: %v", err)
		}
		unstableExps = append(unstableExps, exps)
	}

	// for _, exps := range stableExps {
	// 	experiments.RunExperiments(exps, executable)
	// }

	for _, exps := range unstableExps {
		experiments.RunExperiments(exps, executable)
	}
}
