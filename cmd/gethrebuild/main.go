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
		paths.Scripts.Checkout,
		paths.Scripts.CompareBinaries,
		paths.Scripts.DiffReport,
		paths.Scripts.StartDocker,
		paths.Scripts.Verify,
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

	bi := config.NewBuildInput(af, tc, env, paths)
	fmt.Println(bi)

	err = utils.StartDocker(paths)
	if err != nil {
		log.Fatal(err)
	}

	err = rebuild.RunDockerBuild(bi)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nRebuilding finished, comparing binaries...\n\n")

	result, err := rebuild.Verify(bi.DockerTag, paths)
	if err != nil {
		log.Fatal(err)
	}

	if result.Status == "mismatch" && pa.Diff {
		rebuild.GenerateDiffReport(bi.DockerTag, paths)
	}
}
