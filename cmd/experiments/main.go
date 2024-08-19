package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/chains-project/geth-rebuild/internal/experiments"
	"github.com/chains-project/geth-rebuild/internal/utils"
)

type ExperimentPaths struct {
	RootDir                      string
	GethDir                      string
	MainProgram                  string
	Executable                   string
	GetUnstableCommits           string
	UnstableCommitsFile          string
}

var (
	// All stable builds since 1.14.2, which is the upgrade from ubuntu bionic to noble in travis.yml
	StableVersions = []string{"1.14.8", "1.14.7", "1.14.6", "1.14.5", "1.14.4", "1.14.3", "1.14.2"}
	Arches         = []utils.Arch{utils.AMD64, utils.I386, utils.ARM5, utils.ARM6, utils.ARM7, utils.ARM64}
	paths          ExperimentPaths
)

func init() {
	// set up some utility paths...
	base, err := utils.GetRootDir("geth-rebuild")
	if err != nil {
		log.Fatalf("error finding root path: %v", err)
	}

	paths = ExperimentPaths{
		RootDir:                      base,
		GethDir:                      filepath.Join(base, "tmp", "go-ethereum"),
		MainProgram:                  filepath.Join(base, "cmd", "gethrebuild"),
		Executable:                   filepath.Join(base, "gethrebuild"),
		GetUnstableCommits:           filepath.Join(base, "internal", "experiments", "scripts", "get_available_unstable.sh"),
		UnstableCommitsFile:          filepath.Join(base, "internal", "experiments", "data", "unstable_versions.json"),
	}

	//utils.ChangePermission([]string{paths.GetUnstableCommits}, 0755)

	// build the executable...
	_, err = utils.RunCommand("go", "build", paths.MainProgram)
	if err != nil {
		log.Fatalf("error creating executable: %v", err)
	}
}

func main() {
	exps, err := experiments.GenerateAllExperiments(utils.Linux, Arches, []string{}, paths.UnstableCommitsFile)
	if err != nil {
		log.Fatalf("error generating experiments: %v", err)
	}

	fmt.Println(exps)

	for _, exp := range exps {
		experiments.RunExperiments(exp, paths.Executable)
	}
}
