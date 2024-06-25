package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	util "github.com/chains-project/geth-rebuild/lib"
)

var BuildArgs = map[string]string{
	"UBUNTU_VERSION": "focal",  // default
	"GO_VERSION":     "1.22.0", // default
	"C_COMPILER":     "",
	"GETH_VERSION":   "",
	"OS_ARCH":        "",
	"GETH_COMMIT":    "",
	"BINARY_URL":     "",
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
	// 0. set and validate parameters
	osArch := os.Args[1]
	gethVersion := os.Args[2]
	// repoURL := "https://github.com/ethereum/go-ethereum.git"
	// branch := "master"
	checkoutVersion := "v" + gethVersion

	// validate input parameters
	if err := util.ValidParams(osArch, gethVersion); err != nil {
		log.Fatal(err)
	}

	// directory paths
	rootDir, err := util.GetRootDir()
	if err != nil {
		log.Fatal(err)
	}
	gethDir := rootDir + "/tmp/go-ethereum"
	travisPath := gethDir + "/.travis.yml"
	dockerPath := rootDir + "/docker/rebuilder"

	// 1. clone geth
	fmt.Printf("\n[CLONING GO ETHEREUM SOURCES]\nos-arch		%s\ngeth version	%s\n\n", osArch, gethVersion)
	//util.RunCommand("rm", "-r", "-f", gethDir)
	//util.RunCommand("git", "clone", "--branch", branch, repoURL, gethDir) // TODO: shallow copy. Decide proper --depth OR use --single-branch
	util.RunCommand(gethDir, "git", "fetch")
	util.RunCommand(gethDir, "git", "checkout", checkoutVersion)

	// 2. get hash commit at version for download url
	fmt.Printf("\n[RETRIEVING BINARY DOWNLOAD URL]\n")
	gethCommit := util.RunCommand(rootDir, "git", "log", "-1", "--format=%H")
	gethCommit = strings.ReplaceAll(gethCommit, "\n", "")
	shortCommit := gethCommit[0:8]
	fmt.Printf("\nCommit:		%s", gethCommit)

	targetPackage := "geth-" + osArch + "-" + gethVersion + "-" + shortCommit
	binaryURL := "https://gethstore.blob.core.windows.net/builds/" + targetPackage + ".tar.gz"
	fmt.Printf("\nURL:		%s\n", binaryURL)

	// 3. retrieve necessary build configurations
	fmt.Printf("\n[RETRIEVING BUILD CONFIGURATIONS]\n")
	cc, buildCmd, packages, err := util.GetBuildConfigs(osArch, travisPath)
	if err != nil {
		log.Fatal(err)
	}
	BuildArgs["GETH_VERSION"] = gethVersion
	BuildArgs["OS_ARCH"] = osArch
	BuildArgs["GETH_COMMIT"] = gethCommit
	BuildArgs["BINARY_URL"] = binaryURL
	BuildArgs["BUILD_CMD"] = buildCmd
	BuildArgs["PACKAGES"] = strings.Join(packages, " ")
	BuildArgs["C_COMPILER"] = cc
	for k, v := range BuildArgs {
		fmt.Println(k + ":	" + v)
	}
	// TODO go version
	// TODO ubuntu distribution

	// 4. start verification in docker container
	fmt.Printf("\n[STARTING DOCKER BUILD]\n")
	util.EnsureDocker()
	util.RunDockerBuild(BuildArgs, dockerPath)
}
