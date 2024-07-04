package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chains-project/geth-rebuild/utils"
	"github.com/chains-project/geth-rebuild/utils/specs"
)

type Spec interface {
	ToMap() map[string]string
	PrintSpec() string
}

// build inputs structs
type BuildInput struct {
	Toolchain specs.ToolchainSpec
	Artifact  specs.ArtifactSpec
	Ubuntu    specs.UbuntuSpec
	DockerTag string
}

var paths utils.Paths = utils.SetUpPaths()

func init() {
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
	if len(os.Args) != 4 {
		fmt.Println("Usage: <os> <arch> <geth version>")
		fmt.Println("Example: linux amd64 1.14.3")
		os.Exit(1)
	}
	ops := os.Args[1]
	arch := os.Args[2]
	version := os.Args[3]
	err := utils.ValidateArgs(ops, arch, version)
	if err != nil {
		log.Fatal(err)
	}

	afs, err := specs.NewArtifactSpec(ops, arch, version, paths)
	if err != nil {
		log.Fatal(err)
	}

	// toolchain info
	tc, err := specs.NewToolchainSpec(afs, paths)
	if err != nil {
		log.Fatal(err)
	}

	ub, err := specs.NewUbuntuSpec(afs)
	if err != nil {
		log.Fatal(err)
	}

	buildInput := BuildInput{
		Artifact:  afs,
		Toolchain: tc,
		Ubuntu:    ub,
		DockerTag: utils.CreateDockerTag(afs.Version, afs.Os, afs.Arch),
	}

	_, err = utils.RunCommand(paths.Scripts.StartDocker)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(buildInput)

	//biMap := buildInput.ToMap() // TODO
	//runDockerBuild(biMap, dockerSpec.BuildTag, paths.Directories.Root)
	//utils.RunCommand(copyBinaries, dockerTag, paths.BinDir) // TODO copy into specific dir
	//utils.RunCommand(compareBinaries, paths.BinDir)
}
