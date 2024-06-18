package main

import (
	"fmt"
	"log"
	"os"

	util "github.com/chains-project/geth-rebuild/lib"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: <os-arch> <geth version>")
		fmt.Println("Example: linux-amd64 1.14.3")
		os.Exit(1)
	}

	// TODO runtime.GOOS affects commands?

	osArch := os.Args[1]
	gethVersion := os.Args[2]
	// repoURL := "https://github.com/ethereum/go-ethereum.git"
	// branch := "master"
	rootDir, err := util.GetRootDir()
	if err != nil {
		log.Fatal(err)
	}
	gethDir := rootDir + "/tmp/go-ethereum"
	travisPath := gethDir + "/.travis.yml"

	checkoutVersion := "v" + gethVersion

	if err := util.ValidParams(osArch, gethVersion); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n[CLONING GO ETHEREUM SOURCES]\nos-arch		%s\ngeth version	%s\n\n", osArch, gethVersion)
	//util.RunCommand("rm", "-r", "-f", gethDir)
	//util.RunCommand("git", "clone", "--branch", branch, repoURL, gethDir) // TODO: shallow copy. Decide proper --depth OR use --single-branch
	util.RunCommand(gethDir, "git", "fetch")
	util.RunCommand(gethDir, "git", "checkout", checkoutVersion)

	fmt.Printf("\n[CONSTRUCTING URL FOR BINARY DOWNLOAD]\n")
	commit := util.RunCommand(rootDir, "git", "log", "-1", "--format=%H")
	shortCommit := commit[0:8]

	fmt.Printf("\nCommit:		%sShort commit: 	%s", commit, shortCommit)

	targetPackage := "geth-" + osArch + "-" + gethVersion + "-" + shortCommit
	url := "https://gethstore.blob.core.windows.net/builds/" + targetPackage + ".tar.gz"
	fmt.Printf("\nURL:		%s\n", url)

	fmt.Printf("\n[RETRIEVING BUILD COMMANDS]\n")

	archId := util.GetArchId(osArch) // todo fix (?)
	fmt.Printf("\nArch Id:	%s\n", archId)

	args := util.GetBuildArgs(travisPath, archId)

}
