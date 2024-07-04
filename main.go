package main

import (
	"log"
	"path/filepath"

	"github.com/chains-project/geth-rebuild/utils"
	"github.com/chains-project/geth-rebuild/utils/build"
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
	ops, arch, version, gethDir, unstableHash := utils.ParseFlags()

	err := utils.ValidateArgs(ops, arch, version)
	if err != nil {
		log.Fatal(err)
	}

	var noClone bool // TODO ugly code
	if gethDir != "" {
		paths.Directories.Geth = gethDir
		noClone = true
	}

	// artifact specification
	afs, err := build.NewArtifactSpec(ops, arch, version, unstableHash, noClone, paths)
	if err != nil {
		log.Fatal(err)
	}

	// toolchain specification
	tc, err := build.NewToolchainSpec(afs, paths)
	if err != nil {
		log.Fatal(err)
	}

	// ubuntu specification
	ub, err := build.NewUbuntuSpec(afs, paths)
	if err != nil {
		log.Fatal(err)
	}

	bi := build.BuildInput{
		Artifact:  afs,
		Toolchain: tc,
		Ubuntu:    ub,
		DockerTag: build.CreateDockerTag(afs.Version, afs.Os, afs.Arch),
	}

	_, err = utils.RunCommand(paths.Scripts.StartDocker)
	if err != nil {
		log.Fatal(err)
	}

	build.RunDockerBuild(bi, paths.Directories.Root)
	utils.RunCommand(paths.Scripts.CopyBinaries, bi.DockerTag, paths.Directories.Bin)
	// TODO organise into functions. Alternatively: put scripts into docker.
	binRef := filepath.Join(paths.Directories.Bin, "geth-reference")
	binRep := filepath.Join(paths.Directories.Bin, "geth-reproduce")
	utils.RunCommand(paths.Scripts.CompareBinaries, binRef, binRep)
}
