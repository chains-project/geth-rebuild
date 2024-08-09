package main

import (
	"log"
	"path/filepath"

	"github.com/chains-project/geth-rebuild/internal/experiments"
	"github.com/chains-project/geth-rebuild/internal/utils"
)

type ExperimentPaths struct {
	RootDir     string
	MainProgram string
	Executable  string
	CommitsFile string
}

var (
	// All stable builds since 1.14.2, which is the upgrade from ubuntu bionic to noble in travis.yml
	StableVersions = []string{"1.14.7", "1.14.6", "1.14.5", "1.14.4", "1.14.3", "1.14.2"}
	Arches         = []utils.Arch{utils.AMD64, utils.A386, utils.ARM5, utils.ARM6, utils.ARM7, utils.ARM64}
	paths          ExperimentPaths
)

func init() {
	base, err := utils.GetRootDir("geth-rebuild")
	if err != nil {
		log.Fatalf("error finding root path: %v", err)
	}

	paths = ExperimentPaths{
		RootDir:     base,
		MainProgram: filepath.Join(base, "cmd", "gethrebuild"),
		Executable:  filepath.Join(base, "gethrebuild"),
		CommitsFile: filepath.Join(base, "internal", "experiments", "data", "20_latest_commits.json")}

	_, err = utils.RunCommand("go", "build", paths.MainProgram)
	if err != nil {
		log.Fatalf("error creating executable: %v", err)
	}

}

func main() {
	exps, err := experiments.GenerateAllExperiments(utils.Linux, Arches, StableVersions, paths.CommitsFile)
	if err != nil {
		log.Fatalf("error generating experiments: %v", err)
	}

	for _, exp := range exps {
		experiments.RunExperiments(exp, paths.Executable)
	}
}
