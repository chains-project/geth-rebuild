package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/chains-project/geth-rebuild/common"
)

type Paths struct {
	RootDir    string
	RebuildDir string
	TmpDir     string
	GethDir    string
	ScriptDir  string
	BinDir     string
	TravisPath string
	DockerPath string
}

type BuildArgs struct {
	OsArch        string
	GethVersion   string
	Commit        string
	ShortCommit   string
	GoVersion     string
	CC            string
	UbuntuVersion string
	Packages      string
	BuildCmd      string
	ElfVersion    string
}

// 	DockerTag  string
//	DockerPath string

func validateArgs() (osArch string, gethVersion string) {
	if len(os.Args) != 3 {
		fmt.Println("Usage: <os-arch> <geth version>")
		fmt.Println("Example: linux-amd64 1.14.3") // TODO should change input params to os arch ?
		os.Exit(1)
	}

	osArch = os.Args[1]
	gethVersion = os.Args[2]

	if err := validArgs(osArch, gethVersion); err != nil {
		log.Fatal(err)
	}
	return
}

func main() {
	// validate and set input args
	osArch, gethVersion := validateArgs()

	buildArgs := BuildArgs{
		OsArch:      osArch,
		GethVersion: gethVersion,
	}

	// set up dirs
	rootDir, err := common.GetBaseDir("geth-rebuild")
	if err != nil {
		log.Fatal(err)
	}

	rebuildDir := filepath.Join(rootDir, "cmd", "rebuild")
	tmpDir := filepath.Join(rebuildDir, "tmp")
	gethDir := filepath.Join(tmpDir, "go-ethereum")
	scriptDir := filepath.Join(rebuildDir, "scripts")
	binDir := filepath.Join(rebuildDir, "bin")
	travisPath := filepath.Join(tmpDir, ".travis.yml")
	dockerPath := filepath.Join(rebuildDir, "Dockerfile")

	paths := Paths{
		RootDir:    rootDir,
		RebuildDir: rebuildDir,
		TmpDir:     tmpDir,
		GethDir:    gethDir,
		TravisPath: travisPath,
		ScriptDir:  scriptDir,
		BinDir:     binDir,
		DockerPath: dockerPath,
	}

	// set up scripts
	cloneGeth := filepath.Join(paths.ScriptDir, "clone.sh")
	checkoutGeth := filepath.Join(paths.ScriptDir, "checkout.sh")
	startDocker := filepath.Join(paths.ScriptDir, "start_docker.sh")
	copyBinaries := filepath.Join(paths.ScriptDir, "copy_bin.sh")
	compareBinaries := filepath.Join(paths.ScriptDir, "compare_bin.sh")

	scripts := []string{
		cloneGeth, checkoutGeth, startDocker, copyBinaries, compareBinaries,
	}
	changePermissions(scripts, "+x")

	// clone geth & checkout at version
	//fmt.Printf("\n[CLONING GO ETHEREUM SOURCES]\nos-arch		%s\ngeth version	%s\n\n", osArch, gethVersion)
	common.RunCommand(cloneGeth, paths.TmpDir)
	common.RunCommand(checkoutGeth, paths.GethDir, buildArgs.GethVersion)

	// retrieve build arguments
	//fmt.Printf("\n[RETRIEVING BUILD CONFIGURATIONS FROM SOURCES]\n")
	// commit info
	gethCommit := getCommit(paths.GethDir)
	buildArgs.Commit = gethCommit
	buildArgs.ShortCommit = gethCommit[0:8]

	cc, buildCmd, packages, err := getBuildConfigs(buildArgs.OsArch, paths.TravisPath)
	if err != nil {
		log.Fatal(err)
	}
	buildArgs.CC = cc
	buildArgs.BuildCmd = buildCmd
	buildArgs.Packages = strings.Join(packages, " ")

	buildArgs.GoVersion = "1.22.0"        // TODO
	buildArgs.UbuntuVersion = "focal"     // TODO
	buildArgs.ElfVersion = "elf64-x86-64" // TODO

	// fmt.Print("\n")
	// for k, v := range dockerArgs {
	// 	fmt.Println(k + ":	" + v)
	// }

	dockerTag := createDockerTag(buildArgs.GethVersion, buildArgs.OsArch)
	dockerPath = filepath.Join(paths.RebuildDir, "/Dockerfile")

	// 5. start verification in docker container
	fmt.Printf("\n[STARTING REBUILD IN DOCKER]\n")
	common.RunCommand(startDocker)

	buildArgsMap := buildArgs.ToMap()
	runDockerBuild(buildArgsMap, dockerTag, dockerPath)
	common.RunCommand(copyBinaries, dockerTag, paths.BinDir) // TODO copy into specific dir
	common.RunCommand(compareBinaries, paths.BinDir)
}
