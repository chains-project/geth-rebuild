package main

import (
	"fmt"
	"log"

	"github.com/chains-project/geth-rebuild/internal/buildspec"
	"github.com/chains-project/geth-rebuild/internal/utils"
)

var paths utils.Paths = utils.SetUpPaths()

func init() {
	// set up scripts
	scripts := []string{
		paths.Scripts.Clone,
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
	pa, err := buildspec.ParseArgs()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("OS: %s\n", pa.GOOS)
	fmt.Printf("Arch: %s\n", pa.GOARCH)
	fmt.Printf("Version: %s\n", pa.GethVersion)
	fmt.Printf("Geth Dir: %s\n", pa.GethDir)
	fmt.Printf("Unstable Commit: %s\n", pa.Unstable)

	err = buildspec.ValidateInput(pa)
	if err != nil {
		log.Fatal(err)
	}


	// Validate inputs
	// err := buildspec.ValidateInput(goOs, goArch, gethVersion)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// var noClone bool // TODO ugly code
	// if gethDir != "" {
	// 	paths.Directories.Geth = gethDir
	// 	noClone = true
	// }

	// // artifact specification
	// af, err := config.NewArtifactSpec(ops, arch, version, unstableHash, noClone, paths)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // toolchain specification
	// tc, err := config.NewToolchainSpec(af, paths)
	// if err != nil {
	// 	log.Fatal(err)
	// }

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
