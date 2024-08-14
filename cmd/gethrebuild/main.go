package main

import (
	"fmt"
	"log"

	config "github.com/chains-project/geth-rebuild/internal/buildconfig"
	"github.com/chains-project/geth-rebuild/internal/rebuild"
	"github.com/chains-project/geth-rebuild/internal/utils"
)

var paths utils.Paths = utils.SetUpPaths()

func init() {
	// chmod scripts
	scripts := []string{
		paths.Scripts.GitCheckout,
		paths.Scripts.CompareBinaries,
		paths.Scripts.GenerateDiffReport,
		paths.Scripts.StartDocker,
		paths.Scripts.GetRebuildResults,
	}
	err := utils.ChangePermission(scripts, 0755) // add execute permissions
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	pa, err := utils.ParseArgs()
	if err != nil {
		log.Fatal(err)
	}

	err = utils.ValidArgs(pa)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nRebuilding geth at version %s for %s %s\n\n", fmt.Sprintf("%s %s", pa.GethVersion, pa.Unstable), pa.OS, pa.Arch)

	// artifact specification
	af, err := config.NewArtifactSpec(pa, paths)
	if err != nil {
		log.Fatal(err)
	}

	// toolchain specification
	tc, err := config.NewToolchainSpec(af, paths)
	if err != nil {
		log.Fatal(err)
	}

	// build environment specification
	env, err := config.NewEnvSpec(af, paths)
	if err != nil {
		log.Fatal(err)
	}

	// ensure docker is running
	err = utils.StartDocker(paths)
	if err != nil {
		log.Fatal(err)
	}

	// gather all build inputs
	buildConfig := config.NewBuildConfig(af, tc, env, paths)
	fmt.Println(buildConfig)

	// rebuild in docker
	err = rebuild.DockerRebuild(buildConfig, paths)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nRebuilding finished, comparing binaries...\n\n")

	// Run containerized "verification" i.e. comparison of binaries
	err = rebuild.RunComparison(buildConfig, paths)
	if err != nil {
		log.Fatal(err)
	}

	result, err := rebuild.ReadRebuildResult()
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve the results as logged to file
	err = rebuild.ProcessLogFile(buildConfig.DockerTag, paths, result)
	if err != nil {
		log.Fatal(err)
	}

	// Optional diffoscope for a mismatch
	if result == rebuild.Mismatch && pa.Diff {
		rebuild.GenerateDiffReport(buildConfig.DockerTag, paths)
	}
}
