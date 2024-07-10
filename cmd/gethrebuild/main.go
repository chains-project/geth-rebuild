package main

import (
	"log"

	"github.com/chains-project/geth-rebuild/internal/buildconfig"
	"github.com/chains-project/geth-rebuild/internal/utils"
)

var paths utils.Paths = utils.SetUpPaths()

func init() {
	// set up scripts
	scripts := []string{
		paths.Scripts.Checkout,
		paths.Scripts.StartDocker,
		paths.Scripts.CopyBinaries,
		paths.Scripts.CompareBinaries,
	}
	err := utils.ChangePermissions(scripts, 0755) // add execute permissions
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


	// // artifact specification
	af, err := buildconfig.NewArtifactSpec(pa, paths)
	if err != nil {
		log.Fatal(err)
	}

	// toolchain specification
	_, err = buildconfig.NewToolchainSpec(af, paths)
	if err != nil {
		log.Fatal(err)
	}

	// // ubuntu specification
	// de, err := config.NewDockerSpec(af, paths)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// bi := config.NewBuildInput(af, tc, de, paths)
	// fmt.Println(bi)

	// err = utils.StartDocker(paths)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = rebuild.RunDockerBuild(bi)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// _, err = rebuild.CompareBinaries(bi.DockerTag, paths)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // TODO organise into functions. Alternatively: put scripts into docker.
	// binRef := filepath.Join(paths.Directories.Bin, "geth-reference")
	// binRep := filepath.Join(paths.Directories.Bin, "geth-reproduce")
	// utils.RunCommand(paths.Scripts.CompareBinaries, binRef, binRep)
}
