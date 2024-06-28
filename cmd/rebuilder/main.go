package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	util "github.com/chains-project/geth-rebuild/lib"
)

var rebuildArgs = map[string]string{
	"UBUNTU_VERSION": "focal",  // default
	"GO_VERSION":     "1.22.0", // default
	"C_COMPILER":     "",
	"GETH_VERSION":   "",
	"OS_ARCH":        "",
	"GETH_COMMIT":    "",
	"REFERENCE_URL":  "",
	"BUILD_CMD":      "",
	"ELF_TARGET":     "elf64-x86-64", // todo fix arch specific.
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: <os-arch> <geth version>")
		fmt.Println("Example: linux-amd64 1.14.3")
		os.Exit(1)
	}

	// TODO how runtime.GOOS affects commands
	// 1. Validate input parameters
	osArch := os.Args[1]
	gethVersion := os.Args[2]
	if err := util.ValidParams(osArch, gethVersion); err != nil {
		log.Fatal(err)
	}

	// 2. Set parameters
	// repoURL := "https://github.com/ethereum/go-ethereum.git"
	// branch := "master"
	checkoutVersion := "v" + gethVersion
	rootDir, err := util.GetRootDir()
	if err != nil {
		log.Fatal(err)
	}
	gethDir := rootDir + "/tmp/go-ethereum"
	travisPath := gethDir + "/.travis.yml"
	dockerPath := rootDir + "/docker/rebuilder"

	// 3. clone geth & checkout at version
	fmt.Printf("\n[CLONING GO ETHEREUM SOURCES]\nos-arch		%s\ngeth version	%s\n\n", osArch, gethVersion)
	//util.RunCommand("rm", "-r", "-f", gethDir)
	//util.RunCommand("git", "clone", "--branch", branch, repoURL, gethDir) // TODO: shallow copy. Decide proper --depth OR use --single-branch
	util.RunCommand(gethDir, "git", "fetch")
	util.RunCommand(gethDir, "git", "checkout", checkoutVersion)

	// 4. retrieve all necessary parameters for rebuilding in docker.
	fmt.Printf("\n[RETRIEVING DOCKER BUILD PARAMETERS]\n")
	
	// commit at version
	gethCommit := util.GetCommit(gethDir)
	rebuildArgs["GETH_COMMIT"] = gethCommit
	
	// url for downloading reference binary
	referenceURL := util.GetDownloadURL(osArch, gethVersion, gethCommit)
	rebuildArgs["REFERENCE_URL"] = referenceURL
	
	// build configurations
	cc, buildCmd, packages, err := util.GetBuildConfigs(osArch, travisPath)
	if err != nil {
		log.Fatal(err)
	}
	
	rebuildArgs["GETH_VERSION"] = gethVersion
	rebuildArgs["OS_ARCH"] = osArch
	rebuildArgs["BUILD_CMD"] = buildCmd
	rebuildArgs["PACKAGES"] = strings.Join(packages, " ")
	rebuildArgs["C_COMPILER"] = cc
	
	// print all docker args
	fmt.Print("\n")
	for k, v := range rebuildArgs {
		fmt.Println(k + ":	" + v)
	}
	// TODO go version
	// TODO ubuntu distribution

	// 5. start verification in docker container
	fmt.Printf("\n[STARTING DOCKER BUILD]\n")
	util.EnsureDocker()
	util.RunDockerBuild(rebuildArgs, dockerPath)
}
