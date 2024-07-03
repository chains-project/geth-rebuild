package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chains-project/geth-rebuild/common"
)


var paths Paths = setUpPaths()

func init() {
	scripts := []string{
		paths.Scripts.Clone,
		paths.Scripts.Checkout,
		paths.Scripts.StartDocker,
		paths.Scripts.CopyBinaries,
		paths.Scripts.CompareBinaries,
	}
	err := common.ChangePermissions(scripts, 0755) // add execute permissions
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: <os> <arch> <geth version>")
		fmt.Println("Example: linux amd64 1.14.3")
		os.Exit(1)
	}
	ops := os.Args[1]
	arch := os.Args[2]
	version := os.Args[3]
	err := validateArgs(ops, arch, version)
	if err != nil {
		log.Fatal(err)
	}

	afs, err := getArtifactSpec(ops, arch, version, paths)
	if err != nil {
		log.Fatal(err)
	}

	// toolchain info
	tc, err := getToolchainSpec(afs, paths)
	if err != nil {
		log.Fatal(err)
	}

	ub, err := getUbuntuSpec(afs)
	if err != nil {
		log.Fatal(err)
	}


	buildInput := BuildInput{
		Artifact: afs,
		Toolchain: tc,
		Ubuntu: ub,
	}

	dockerSpec := DockerSpec{
		Dir: paths.Files.Docker,
		FileHash: "", // TODO
		BuildTag:  createDockerTag(buildInput.Artifact.Version, buildInput.Artifact.Os, buildInput.Artifact.Arch),
	}


	common.RunCommand(paths.Scripts.StartDocker)
	
	biMap := buildInput.ToMap()
	runDockerBuild(biMap, dockerSpec.BuildTag, paths.Directories.Rebuild)
	//common.RunCommand(copyBinaries, dockerTag, paths.BinDir) // TODO copy into specific dir
	//common.RunCommand(compareBinaries, paths.BinDir)
}
