package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/chains-project/geth-rebuild/internal/rebuild"
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
	afs, err := rebuild.NewArtifactSpec(ops, arch, version, unstableHash, noClone, paths)
	if err != nil {
		log.Fatal(err)
	}

	// toolchain specification
	tc, err := rebuild.NewToolchainSpec(afs, paths)
	if err != nil {
		log.Fatal(err)
	}

	// ubuntu specification
	ub, err := rebuild.NewDockerSpec(afs, paths)
	if err != nil {
		log.Fatal(err)
	}

	bi := rebuild.BuildInput{
		Artifact:  afs,
		Toolchain: tc,
		Ubuntu:    ub,
		DockerTag: rebuild.CreateDockerTag(afs.Version, afs.Os, afs.Arch),
	}

	fmt.Println(bi)

	_, err = utils.RunCommand(paths.Scripts.StartDocker)
	if err != nil {
		log.Fatal(err)
	}

	bi.RunDockerBuild(paths.Directories.Root)
	utils.RunCommand(paths.Scripts.CopyBinaries, bi.DockerTag, paths.Directories.Bin)
	// TODO organise into functions. Alternatively: put scripts into docker.
	binRef := filepath.Join(paths.Directories.Bin, "geth-reference")
	binRep := filepath.Join(paths.Directories.Bin, "geth-reproduce")
	utils.RunCommand(paths.Scripts.CompareBinaries, binRef, binRep)
}
